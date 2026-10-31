package impl

import (
	"context"
	"reflect"
	"testing"

	"ads-platform/internal/business/iam/staff/model"
	"ads-platform/internal/business/iam/staff/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/mobile"

	"gorm.io/gorm"
)

type fakeRepo struct {
	next  int64
	staff map[int64]*model.Staff
}

func (f *fakeRepo) CreateStaff(_ context.Context, staff *model.Staff) error {
	f.next++
	staff.ID = f.next
	cp := *staff
	f.staff[staff.ID] = &cp
	return nil
}

func (f *fakeRepo) GetStaffByID(_ context.Context, id int64) (*model.Staff, error) {
	s, ok := f.staff[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *s
	return &cp, nil
}

func (f *fakeRepo) ListStaff(context.Context, []string, int, int) ([]model.Staff, int64, error) {
	return nil, 0, nil
}
func (f *fakeRepo) UpdateStaff(_ context.Context, staff *model.Staff) error {
	cp := *staff
	f.staff[staff.ID] = &cp
	return nil
}

func (f *fakeRepo) DeleteStaff(_ context.Context, id int64) error {
	if _, ok := f.staff[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(f.staff, id)
	return nil
}

func (f *fakeRepo) GetStaffByMobile(_ context.Context, mobile string) (*model.Staff, error) {
	for _, s := range f.staff {
		if s.Mobile == mobile {
			cp := *s
			return &cp, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func newService() service.StaffService {
	return newServiceWith(&fakeRepo{staff: map[int64]*model.Staff{}})
}

func newServiceWith(repo *fakeRepo) service.StaffService {
	return NewStaffService(repo, mobile.NewNormalizer("+98"))
}

func errorCode(t *testing.T, err error) string {
	t.Helper()
	appErr, ok := exception.AsAppError(err)
	if !ok {
		t.Fatalf("want AppError, got %v", err)
	}
	return appErr.ErrorCode
}

func TestCreateNormalizesInput(t *testing.T) {
	got, err := newService().Create(context.Background(), service.StaffInput{
		Name:       "  علی رضایی ",
		Mobile:     "۰۹۱۲۱۲۳۴۵۶۷",
		NationalId: "۰۰۱۲۳۴۵۶۷۸",
		Password:   "a-strong-pass",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "علی رضایی" || got.Mobile != "+989121234567" || got.NationalId != "0012345678" {
		t.Fatalf("got %+v", got)
	}
	if !got.IsActive || got.IsSuperAdmin || got.PasswordHash == "" || got.PasswordHash == "a-strong-pass" {
		t.Fatalf("unexpected auth fields %+v", got)
	}
}

func TestCreateRejectsInvalidInput(t *testing.T) {
	cases := map[string]service.StaffInput{
		"STAFF_INVALID_NAME":        {Name: " ", Mobile: "09121234567", NationalId: "0012345678", Password: "a-strong-pass"},
		"STAFF_INVALID_MOBILE":      {Name: "Ali", Mobile: "12", NationalId: "0012345678", Password: "a-strong-pass"},
		"STAFF_INVALID_NATIONAL_ID": {Name: "Ali", Mobile: "09121234567", NationalId: "123", Password: "a-strong-pass"},
		"STAFF_INVALID_PASSWORD":    {Name: "Ali", Mobile: "09121234567", NationalId: "0012345678", Password: "short"},
	}
	for want, input := range cases {
		_, err := newService().Create(context.Background(), input)
		if got := errorCode(t, err); got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	}
}

func TestDeleteMissingStaff(t *testing.T) {
	err := newService().Delete(context.Background(), 42)
	if got := errorCode(t, err); got != "STAFF_NOT_FOUND" {
		t.Fatalf("got %s", got)
	}
}

func TestSearchTermsMatchLocalMobile(t *testing.T) {
	if got := searchTerms("۰۹۱۲"); !reflect.DeepEqual(got, []string{"0912", "912"}) {
		t.Fatalf("got %v", got)
	}
	if got := searchTerms(" علی "); !reflect.DeepEqual(got, []string{"علی"}) {
		t.Fatalf("got %v", got)
	}
	if got := searchTerms(""); got != nil {
		t.Fatalf("got %v", got)
	}
}

func TestSuperAdminIsProtected(t *testing.T) {
	repo := &fakeRepo{staff: map[int64]*model.Staff{}}
	svc := newServiceWith(repo)
	admin, err := svc.UpsertSuperAdmin(context.Background(), service.StaffInput{
		Name: "Admin", Mobile: "09120000001", NationalId: "0000000001", Password: "a-strong-pass",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !admin.IsSuperAdmin || !admin.IsActive {
		t.Fatalf("got %+v", admin)
	}

	if got := errorCode(t, svc.Delete(context.Background(), admin.ID)); got != "STAFF_SUPER_ADMIN_PROTECTED" {
		t.Fatalf("delete: got %s", got)
	}
	inactive := false
	_, err = svc.Update(context.Background(), admin.ID, service.StaffInput{
		Name: "Admin", Mobile: "09120000001", NationalId: "0000000001", IsActive: &inactive,
	})
	if got := errorCode(t, err); got != "STAFF_SUPER_ADMIN_PROTECTED" {
		t.Fatalf("deactivate: got %s", got)
	}
}

func TestUpsertSuperAdminPromotesExistingStaff(t *testing.T) {
	repo := &fakeRepo{staff: map[int64]*model.Staff{}}
	svc := newServiceWith(repo)
	input := service.StaffInput{Name: "Sara", Mobile: "09120000002", NationalId: "0000000002", Password: "first-password"}
	staff, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	oldHash := staff.PasswordHash

	input.Password = "second-password"
	promoted, err := svc.UpsertSuperAdmin(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if promoted.ID != staff.ID || !promoted.IsSuperAdmin || promoted.PasswordHash == oldHash {
		t.Fatalf("got %+v", promoted)
	}
}
