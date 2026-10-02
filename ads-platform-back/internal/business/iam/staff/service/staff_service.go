package service

import (
	"ads-platform/internal/business/iam/staff/model"
	"context"
)

type StaffService interface {
	GetStaffByID(ctx context.Context, id int64) (*model.Staff, error)
}
