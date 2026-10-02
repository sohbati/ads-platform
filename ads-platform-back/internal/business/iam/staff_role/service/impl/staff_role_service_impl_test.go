package impl

import (
	"context"
	"testing"

	actionmodel "ads-platform/internal/business/iam/action/model"
	rolemodel "ads-platform/internal/business/iam/role/model"
	staffmodel "ads-platform/internal/business/iam/staff/model"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type fakeStaff struct {
	ids map[int64]struct{}
}

func (f *fakeStaff) GetStaffByID(_ context.Context, id int64) (*staffmodel.Staff, error) {
	if _, ok := f.ids[id]; !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &staffmodel.Staff{ID: id}, nil
}

type fakeRoles struct {
	ids map[int64]struct{}
}

func (f *fakeRoles) GetByID(_ context.Context, id int64) (*rolemodel.Role, error) {
	if _, ok := f.ids[id]; !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return &rolemodel.Role{ID: id, Code: "moderator"}, nil
}

type fakeRepo struct {
	roles    map[int64]map[int64]struct{}
	actions  []actionmodel.Action
	hasCodes map[string]bool
}

func (f *fakeRepo) ListRoles(_ context.Context, staffID int64) ([]rolemodel.Role, error) {
	out := []rolemodel.Role{}
	for roleID := range f.roles[staffID] {
		out = append(out, rolemodel.Role{ID: roleID})
	}
	return out, nil
}

func (f *fakeRepo) Add(_ context.Context, staffID, roleID int64) error {
	if f.roles[staffID] == nil {
		f.roles[staffID] = map[int64]struct{}{}
	}
	f.roles[staffID][roleID] = struct{}{}
	return nil
}

func (f *fakeRepo) Remove(context.Context, int64, int64) error { return nil }

func (f *fakeRepo) ListActions(context.Context, int64) ([]actionmodel.Action, error) {
	if f.actions == nil {
		return []actionmodel.Action{}, nil
	}
	return f.actions, nil
}

func (f *fakeRepo) HasAction(_ context.Context, _ int64, code string) (bool, error) {
	return f.hasCodes[code], nil
}

func TestUnknownStaff(t *testing.T) {
	svc := NewStaffRoleService(
		&fakeRepo{roles: map[int64]map[int64]struct{}{}},
		&fakeStaff{ids: map[int64]struct{}{1: {}}},
		&fakeRoles{ids: map[int64]struct{}{2: {}}},
	)
	_, err := svc.ListActions(context.Background(), 99)
	app, ok := exception.AsAppError(err)
	if !ok || app.ErrorCode != "STAFF_NOT_FOUND" {
		t.Fatalf("err=%v", err)
	}
}

func TestRequireActionUsesAssignedRolesOnly(t *testing.T) {
	svc := NewStaffRoleService(
		&fakeRepo{
			roles:    map[int64]map[int64]struct{}{},
			hasCodes: map[string]bool{"ad.approve": true},
		},
		&fakeStaff{ids: map[int64]struct{}{1: {}}},
		&fakeRoles{ids: map[int64]struct{}{2: {}}},
	)
	if err := svc.RequireAction(context.Background(), 1, "ad.approve"); err != nil {
		t.Fatal(err)
	}
	err := svc.RequireAction(context.Background(), 1, "iam.manage")
	app, ok := exception.AsAppError(err)
	if !ok || app.ErrorCode != "ACTION_DENIED" {
		t.Fatalf("err=%v", err)
	}
}
