package repository

import (
	"ads-platform/internal/business/iam/staff/model"
	"context"
)

type StaffRepository interface {
	CreateStaff(ctx context.Context, staff *model.Staff) error
	GetStaffByID(ctx context.Context, id int64) (*model.Staff, error)
	GetStaffList(ctx context.Context, limit, offset int) ([]*model.Staff, error)
	UpdateStaff(ctx context.Context, staff *model.Staff) error
	DeleteStaff(ctx context.Context, id int64) error
	GetStaffByMobile(ctx context.Context, mobile string) (*model.Staff, error)
}
