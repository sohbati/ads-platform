package repository

import (
	"context"

	"ads-platform/internal/business/iam/staff_group/model"
)

type StaffGroupRepository interface {
	Create(ctx context.Context, group *model.StaffGroup) error
	GetByID(ctx context.Context, id int64) (*model.StaffGroup, error)
	List(ctx context.Context) ([]model.StaffGroup, error)
	Update(ctx context.Context, group *model.StaffGroup) error
	Delete(ctx context.Context, id int64) error
	ListStaffIDs(ctx context.Context, groupID int64) ([]int64, error)
	AddStaff(ctx context.Context, groupID, staffID int64) error
	RemoveStaff(ctx context.Context, groupID, staffID int64) error
	ListByStaff(ctx context.Context, staffID int64) ([]model.StaffGroup, error)
}
