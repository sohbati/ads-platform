package repository

import (
	"ads-platform/internal/business/iam/staff/model"
	"context"
)

type StaffRepository interface {
	CreateStaff(ctx context.Context, staff *model.Staff) error
	GetStaffByID(ctx context.Context, id int64) (*model.Staff, error)
	// ListStaff matches search against name, mobile, and national id; an empty search returns everyone.
	ListStaff(ctx context.Context, search []string, limit, offset int) ([]model.Staff, int64, error)
	UpdateStaff(ctx context.Context, staff *model.Staff) error
	DeleteStaff(ctx context.Context, id int64) error
	GetStaffByMobile(ctx context.Context, mobile string) (*model.Staff, error)
}
