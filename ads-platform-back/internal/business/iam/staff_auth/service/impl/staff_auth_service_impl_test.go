package impl

import (
	"context"
	"reflect"
	"testing"

	actionmodel "ads-platform/internal/business/iam/action/model"
	"ads-platform/internal/business/iam/internal/password"
	staffmodel "ads-platform/internal/business/iam/staff/model"
	"ads-platform/internal/business/iam/staff_auth/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/mobile"

	"gorm.io/gorm"
)

type fakeStaff map[int64]*staffmodel.Staff

func (f fakeStaff) GetStaffByID(_ context.Context, id int64) (*staffmodel.Staff, error) {
	if s, ok := f[id]; ok {
		cp := *s
		return &cp, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f fakeStaff) GetStaffByMobile(_ context.Context, m string) (*staffmodel.Staff, error) {
	for _, s := range f {
		if s.Mobile == m {
			cp := *s
			return &cp, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

type fakeActions map[int64][]string

func (f fakeActions) ListActions(_ context.Context, staffID int64) ([]actionmodel.Action, error) {
	items := []actionmodel.Action{}
	for _, code := range f[staffID] {
		items = append(items, actionmodel.Action{Code: code})
	}
	return items, nil
}

func newAuth(t *testing.T) service.StaffAuthService {
	t.Helper()
	hash, err := password.Hash("a-strong-pass")
	if err != nil {
		t.Fatal(err)
	}
	staff := fakeStaff{
		1: {ID: 1, Mobile: "+989120000001", PasswordHash: hash, IsSuperAdmin: true, IsActive: true},
		2: {ID: 2, Mobile: "+989120000002", PasswordHash: hash, IsActive: true},
		3: {ID: 3, Mobile: "+989120000003", PasswordHash: hash, IsActive: false},
		4: {ID: 4, Mobile: "+989120000004", IsActive: true},
	}
	actions := fakeActions{2: {"staff.read", "role.read"}}
	return NewStaffAuthService(staff, actions, mobile.NewNormalizer("+98"))
}

func code(t *testing.T, err error) string {
	t.Helper()
	appErr, ok := exception.AsAppError(err)
	if !ok {
		t.Fatalf("want AppError, got %v", err)
	}
	return appErr.ErrorCode
}

func TestLoginSuperAdminGetsWildcard(t *testing.T) {
	session, err := newAuth(t).Login(context.Background(), "۰۹۱۲۰۰۰۰۰۰۱", "a-strong-pass")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(session.Permissions, []string{"*"}) {
		t.Fatalf("got %v", session.Permissions)
	}
}

func TestLoginStaffGetsRoleActions(t *testing.T) {
	session, err := newAuth(t).Login(context.Background(), "09120000002", "a-strong-pass")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(session.Permissions, []string{"staff.read", "role.read"}) {
		t.Fatalf("got %v", session.Permissions)
	}
}

func TestLoginRejects(t *testing.T) {
	cases := map[string][2]string{
		"wrong password": {"09120000002", "wrong-password"},
		"unknown mobile": {"09129999999", "a-strong-pass"},
		"invalid mobile": {"12", "a-strong-pass"},
		"inactive staff": {"09120000003", "a-strong-pass"},
		"no password":    {"09120000004", ""},
	}
	auth := newAuth(t)
	for name, c := range cases {
		_, err := auth.Login(context.Background(), c[0], c[1])
		if got := code(t, err); got != "STAFF_AUTH_INVALID_CREDENTIALS" {
			t.Errorf("%s: got %s", name, got)
		}
	}
}

func TestSessionRejectsInactiveAndMissing(t *testing.T) {
	auth := newAuth(t)
	for _, id := range []int64{0, 3, 99} {
		_, err := auth.Session(context.Background(), id)
		if got := code(t, err); got != "STAFF_AUTH_SESSION_INVALID" {
			t.Errorf("id %d: got %s", id, got)
		}
	}
	if _, err := auth.Session(context.Background(), 2); err != nil {
		t.Fatal(err)
	}
}
