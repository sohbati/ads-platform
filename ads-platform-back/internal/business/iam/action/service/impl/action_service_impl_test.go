package impl

import (
	"context"
	"testing"

	"ads-platform/internal/business/iam/action/model"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type fakeRepo struct {
	next    int64
	actions map[int64]*model.Action
}

func (f *fakeRepo) Create(_ context.Context, action *model.Action) error {
	f.next++
	action.ID = f.next
	cp := *action
	f.actions[action.ID] = &cp
	return nil
}

func (f *fakeRepo) GetByID(_ context.Context, id int64) (*model.Action, error) {
	a, ok := f.actions[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	cp := *a
	return &cp, nil
}

func (f *fakeRepo) List(context.Context) ([]model.Action, error) { return nil, nil }
func (f *fakeRepo) Update(context.Context, *model.Action) error  { return nil }
func (f *fakeRepo) Delete(context.Context, int64) error          { return nil }

func TestCreateNormalizesCode(t *testing.T) {
	svc := NewActionService(&fakeRepo{actions: map[int64]*model.Action{}})
	got, err := svc.Create(context.Background(), " Ad.Approve ", "Approve ad")
	if err != nil {
		t.Fatal(err)
	}
	if got.Code != "ad.approve" {
		t.Fatalf("code=%s", got.Code)
	}
}

func TestCreateRejectsBadCode(t *testing.T) {
	svc := NewActionService(&fakeRepo{actions: map[int64]*model.Action{}})
	_, err := svc.Create(context.Background(), "Approve Ad", "Approve")
	app, ok := exception.AsAppError(err)
	if !ok || app.ErrorCode != "ACTION_INVALID_CODE" {
		t.Fatalf("err=%v", err)
	}
}
