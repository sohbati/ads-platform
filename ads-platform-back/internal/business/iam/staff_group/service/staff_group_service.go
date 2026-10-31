package service

import (
	"context"

	"ads-platform/internal/business/iam/staff_group/model"
)

type StaffGroupService interface {
	Create(ctx context.Context, code, name string) (*model.StaffGroup, error)
	Get(ctx context.Context, id int64) (*model.StaffGroup, error)
	List(ctx context.Context) ([]model.StaffGroup, error)
	// Update renames the item and, when code is non-empty, changes its code too.
	Update(ctx context.Context, id int64, code, name string) (*model.StaffGroup, error)
	Delete(ctx context.Context, id int64) error
	ListStaff(ctx context.Context, groupID int64) ([]int64, error)
	AddStaff(ctx context.Context, groupID, staffID int64) error
	RemoveStaff(ctx context.Context, groupID, staffID int64) error
	ListByStaff(ctx context.Context, staffID int64) ([]model.StaffGroup, error)
}
