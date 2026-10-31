package impl

import (
	"context"
	"testing"

	"ads-platform/internal/business/iam/role/model"
	"ads-platform/internal/core/exception"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type fakeRepo struct {
	roles map[int64]*model.Role
}

func (f *fakeRepo) Create(context.Context, *model.Role) error { return nil }

func (f *fakeRepo) GetByID(_ context.Context, id int64) (*model.Role, error) {
	r, ok := f.roles[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *r
	return &cp, nil
}

func (f *fakeRepo) List(context.Context) ([]model.Role, error) { return nil, nil }

func (f *fakeRepo) Update(_ context.Context, role *model.Role) error {
	for id, other := range f.roles {
		if id != role.ID && other.Code == role.Code {
			return &pgconn.PgError{Code: "23505"}
		}
	}
	cp := *role
	f.roles[role.ID] = &cp
	return nil
}

func (f *fakeRepo) Delete(context.Context, int64) error { return nil }

func newRepo() *fakeRepo {
	return &fakeRepo{roles: map[int64]*model.Role{
		1: {ID: 1, Code: "moderator", Name: "Moderator"},
		2: {ID: 2, Code: "support", Name: "Support"},
	}}
}

func errorCode(t *testing.T, err error) string {
	t.Helper()
	app, ok := exception.AsAppError(err)
	if !ok {
		t.Fatalf("want AppError, got %v", err)
	}
	return app.ErrorCode
}

func TestUpdateChangesCode(t *testing.T) {
	got, err := NewRoleService(newRepo()).Update(context.Background(), 1, " Content.Moderator ", "Content moderator")
	if err != nil {
		t.Fatal(err)
	}
	if got.Code != "content.moderator" || got.Name != "Content moderator" {
		t.Fatalf("got %+v", got)
	}
}

func TestUpdateWithoutCodeKeepsIt(t *testing.T) {
	got, err := NewRoleService(newRepo()).Update(context.Background(), 1, "", "Renamed")
	if err != nil {
		t.Fatal(err)
	}
	if got.Code != "moderator" || got.Name != "Renamed" {
		t.Fatalf("got %+v", got)
	}
}

func TestUpdateRejectsBadOrDuplicateCode(t *testing.T) {
	svc := NewRoleService(newRepo())
	if _, err := svc.Update(context.Background(), 1, "9bad code", "x y"); errorCode(t, err) != "ROLE_INVALID_CODE" {
		t.Fatalf("bad code: %v", err)
	}
	if _, err := svc.Update(context.Background(), 1, "support", "x y"); errorCode(t, err) != "ROLE_DUPLICATE_CODE" {
		t.Fatalf("duplicate: %v", err)
	}
}
