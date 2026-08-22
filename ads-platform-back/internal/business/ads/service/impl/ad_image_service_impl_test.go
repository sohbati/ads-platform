package impl

import (
	"context"
	"strings"
	"testing"

	"ads-platform/internal/business/ads/model"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

// fakeRepo is an in-memory AdImageRepository.
type fakeRepo struct {
	nextID int64
	rows   map[int64]model.AdImage
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{nextID: 1, rows: map[int64]model.AdImage{}}
}

func (f *fakeRepo) Create(_ context.Context, image *model.AdImage) error {
	image.ID = f.nextID
	f.nextID++
	f.rows[image.ID] = *image
	return nil
}

func (f *fakeRepo) GetByID(_ context.Context, id int64) (*model.AdImage, error) {
	row, ok := f.rows[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copy := row
	return &copy, nil
}

func (f *fakeRepo) Update(_ context.Context, image *model.AdImage) error {
	f.rows[image.ID] = *image
	return nil
}

func assertAppErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %s, got nil", code)
	}
	appErr, ok := exception.AsAppError(err)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.ErrorCode != code {
		t.Fatalf("expected error code %s, got %s", code, appErr.ErrorCode)
	}
}

func TestRegisterCreatesPendingImage(t *testing.T) {
	svc := NewAdImageService(newFakeRepo())

	image, err := svc.Register(context.Background(), 42, "my photo.JPG", "image/jpeg", 1024)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if image.ID == 0 {
		t.Error("expected assigned ID")
	}
	if image.Status != model.ImageStatusPending {
		t.Errorf("Status = %q, want pending", image.Status)
	}
	if image.UserID != 42 || image.OriginalFilename != "my photo.JPG" || image.FileSize != 1024 {
		t.Errorf("metadata not stored: %+v", image)
	}
	if !strings.HasPrefix(image.ObjectKey, "ads/42/") || !strings.HasSuffix(image.ObjectKey, ".jpg") {
		t.Errorf("ObjectKey = %q, want ads/42/<random>.jpg", image.ObjectKey)
	}
	if image.Checksum != nil || image.UploadedAt != nil || image.DeletedAt != nil {
		t.Errorf("upload/delete fields must start empty: %+v", image)
	}
}

func TestRegisterObjectKeysAreUnique(t *testing.T) {
	svc := NewAdImageService(newFakeRepo())

	a, _ := svc.Register(context.Background(), 1, "a.png", "image/png", 10)
	b, _ := svc.Register(context.Background(), 1, "a.png", "image/png", 10)
	if a.ObjectKey == b.ObjectKey {
		t.Fatalf("object keys must differ, both %q", a.ObjectKey)
	}
}

func TestRegisterValidation(t *testing.T) {
	svc := NewAdImageService(newFakeRepo())
	ctx := context.Background()

	cases := []struct {
		name        string
		userID      int64
		filename    string
		contentType string
		fileSize    int64
		wantCode    string
	}{
		{"missing user", 0, "a.jpg", "image/jpeg", 10, "IMAGE_INVALID_USER"},
		{"empty filename", 1, "  ", "image/jpeg", 10, "IMAGE_INVALID_FILENAME"},
		{"filename too long", 1, strings.Repeat("x", 256), "image/jpeg", 10, "IMAGE_INVALID_FILENAME"},
		{"disallowed type", 1, "a.exe", "application/octet-stream", 10, "IMAGE_INVALID_CONTENT_TYPE"},
		{"svg rejected", 1, "a.svg", "image/svg+xml", 10, "IMAGE_INVALID_CONTENT_TYPE"},
		{"zero size", 1, "a.jpg", "image/jpeg", 0, "IMAGE_INVALID_SIZE"},
		{"too large", 1, "a.jpg", "image/jpeg", MaxImageSize + 1, "IMAGE_TOO_LARGE"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Register(ctx, tc.userID, tc.filename, tc.contentType, tc.fileSize)
			assertAppErrorCode(t, err, tc.wantCode)
		})
	}
}

func TestMarkUploadedTransition(t *testing.T) {
	repo := newFakeRepo()
	svc := NewAdImageService(repo)
	ctx := context.Background()

	image, _ := svc.Register(ctx, 7, "a.webp", "image/webp", 100)

	updated, err := svc.MarkUploaded(ctx, 7, image.ID, "abc123", 2048)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Status != model.ImageStatusUploaded {
		t.Errorf("Status = %q, want uploaded", updated.Status)
	}
	if updated.Checksum == nil || *updated.Checksum != "abc123" {
		t.Errorf("Checksum = %v, want abc123", updated.Checksum)
	}
	if updated.UploadedAt == nil {
		t.Error("UploadedAt not set")
	}
	if updated.FileSize != 2048 {
		t.Errorf("FileSize = %d, want 2048 (actual uploaded size)", updated.FileSize)
	}

	// Second confirmation must fail: image is no longer pending.
	_, err = svc.MarkUploaded(ctx, 7, image.ID, "abc123", 2048)
	assertAppErrorCode(t, err, "IMAGE_INVALID_STATUS")
}

func TestMarkUploadedRequiresChecksum(t *testing.T) {
	svc := NewAdImageService(newFakeRepo())
	_, err := svc.MarkUploaded(context.Background(), 7, 1, "  ", 0)
	assertAppErrorCode(t, err, "IMAGE_CHECKSUM_REQUIRED")
}

func TestOwnershipIsEnforced(t *testing.T) {
	repo := newFakeRepo()
	svc := NewAdImageService(repo)
	ctx := context.Background()

	image, _ := svc.Register(ctx, 7, "a.jpg", "image/jpeg", 100)

	// Another user sees not-found, not forbidden, so ids can't be probed.
	_, err := svc.Get(ctx, 8, image.ID)
	assertAppErrorCode(t, err, "IMAGE_NOT_FOUND")
	_, err = svc.MarkUploaded(ctx, 8, image.ID, "abc", 0)
	assertAppErrorCode(t, err, "IMAGE_NOT_FOUND")
	err = svc.Delete(ctx, 8, image.ID)
	assertAppErrorCode(t, err, "IMAGE_NOT_FOUND")

	// Unknown id behaves the same.
	_, err = svc.Get(ctx, 7, 9999)
	assertAppErrorCode(t, err, "IMAGE_NOT_FOUND")
}

func TestDeleteIsSoftAndIdempotent(t *testing.T) {
	repo := newFakeRepo()
	svc := NewAdImageService(repo)
	ctx := context.Background()

	image, _ := svc.Register(ctx, 7, "a.jpg", "image/jpeg", 100)

	if err := svc.Delete(ctx, 7, image.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	row := repo.rows[image.ID]
	if row.Status != model.ImageStatusDeleted || row.DeletedAt == nil {
		t.Errorf("expected soft delete, got status %q deleted_at %v", row.Status, row.DeletedAt)
	}

	// Repeating the delete is a no-op, not an error.
	if err := svc.Delete(ctx, 7, image.ID); err != nil {
		t.Fatalf("second delete: unexpected error: %v", err)
	}
}
