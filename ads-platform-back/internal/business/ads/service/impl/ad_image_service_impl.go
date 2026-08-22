package impl

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"ads-platform/internal/business/ads/errorcode"
	"ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/ads/repository"
	"ads-platform/internal/business/ads/service"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

// MaxImageSize is the largest accepted image file (bytes).
const MaxImageSize = 10 << 20 // 10 MB

// extByContentType doubles as the allowlist of accepted image types.
var extByContentType = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type adImageService struct {
	repo repository.AdImageRepository
}

func NewAdImageService(repo repository.AdImageRepository) service.AdImageService {
	return &adImageService{repo: repo}
}

func (s *adImageService) Register(ctx context.Context, userID int64, originalFilename, contentType string, fileSize int64) (*model.AdImage, error) {
	if userID <= 0 {
		return nil, exception.NewAppError(
			errorcode.ErrImageInvalidUser.Code, errorcode.ErrImageInvalidUser.HttpStatus)
	}

	originalFilename = strings.TrimSpace(originalFilename)
	if originalFilename == "" || len(originalFilename) > 255 {
		return nil, exception.NewAppError(
			errorcode.ErrImageInvalidFilename.Code, errorcode.ErrImageInvalidFilename.HttpStatus, originalFilename)
	}

	contentType = strings.ToLower(strings.TrimSpace(contentType))
	ext, ok := extByContentType[contentType]
	if !ok {
		return nil, exception.NewAppError(
			errorcode.ErrImageInvalidType.Code, errorcode.ErrImageInvalidType.HttpStatus, contentType)
	}

	if fileSize <= 0 {
		return nil, exception.NewAppError(
			errorcode.ErrImageInvalidSize.Code, errorcode.ErrImageInvalidSize.HttpStatus)
	}
	if fileSize > MaxImageSize {
		return nil, exception.NewAppError(
			errorcode.ErrImageTooLarge.Code, errorcode.ErrImageTooLarge.HttpStatus)
	}

	objectKey, err := buildObjectKey(userID, ext)
	if err != nil {
		return nil, err
	}

	image := &model.AdImage{
		UserID:           userID,
		ObjectKey:        objectKey,
		OriginalFilename: originalFilename,
		ContentType:      contentType,
		FileSize:         fileSize,
		Status:           model.ImageStatusPending,
	}
	if err := s.repo.Create(ctx, image); err != nil {
		return nil, err
	}
	return image, nil
}

func (s *adImageService) MarkUploaded(ctx context.Context, userID, imageID int64, checksum string, fileSize int64) (*model.AdImage, error) {
	checksum = strings.TrimSpace(checksum)
	if checksum == "" || len(checksum) > 64 {
		return nil, exception.NewAppError(
			errorcode.ErrImageChecksumRequired.Code, errorcode.ErrImageChecksumRequired.HttpStatus)
	}

	image, err := s.ownedImage(ctx, userID, imageID)
	if err != nil {
		return nil, err
	}
	if image.Status != model.ImageStatusPending {
		return nil, exception.NewAppError(
			errorcode.ErrImageInvalidStatus.Code, errorcode.ErrImageInvalidStatus.HttpStatus, image.Status)
	}
	if fileSize > 0 {
		if fileSize > MaxImageSize {
			return nil, exception.NewAppError(
				errorcode.ErrImageTooLarge.Code, errorcode.ErrImageTooLarge.HttpStatus)
		}
		image.FileSize = fileSize
	}

	now := time.Now().UTC()
	image.Status = model.ImageStatusUploaded
	image.Checksum = &checksum
	image.UploadedAt = &now

	if err := s.repo.Update(ctx, image); err != nil {
		return nil, err
	}
	return image, nil
}

func (s *adImageService) Get(ctx context.Context, userID, imageID int64) (*model.AdImage, error) {
	return s.ownedImage(ctx, userID, imageID)
}

func (s *adImageService) Delete(ctx context.Context, userID, imageID int64) error {
	image, err := s.ownedImage(ctx, userID, imageID)
	if err != nil {
		return err
	}
	if image.Status == model.ImageStatusDeleted {
		return nil // idempotent
	}

	now := time.Now().UTC()
	image.Status = model.ImageStatusDeleted
	image.DeletedAt = &now
	return s.repo.Update(ctx, image)
}

// ownedImage loads an image and hides other users' records behind not-found,
// so image ids cannot be probed across accounts.
func (s *adImageService) ownedImage(ctx context.Context, userID, imageID int64) (*model.AdImage, error) {
	if userID <= 0 {
		return nil, exception.NewAppError(
			errorcode.ErrImageInvalidUser.Code, errorcode.ErrImageInvalidUser.HttpStatus)
	}

	image, err := s.repo.GetByID(ctx, imageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(
				errorcode.ErrImageNotFound.Code, errorcode.ErrImageNotFound.HttpStatus).WithCause(err)
		}
		return nil, err
	}
	if image.UserID != userID {
		return nil, exception.NewAppError(
			errorcode.ErrImageNotFound.Code, errorcode.ErrImageNotFound.HttpStatus)
	}
	return image, nil
}

// buildObjectKey creates a storage key like "ads/42/1a2b...cd.jpg"; the random
// component makes keys unguessable and collision-free.
func buildObjectKey(userID int64, ext string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("ad image: generate object key: %w", err)
	}
	return fmt.Sprintf("ads/%d/%s%s", userID, hex.EncodeToString(buf), ext), nil
}
