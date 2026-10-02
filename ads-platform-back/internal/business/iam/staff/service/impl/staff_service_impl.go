package impl

import (
	"ads-platform/internal/business/iam/staff/model"
	"ads-platform/internal/business/iam/staff/repository"
	"ads-platform/internal/business/iam/staff/service"
	"context"
	"errors"
)

type staffService struct {
	staffRepo repository.StaffRepository
}

func NewStaffService(staffRepo repository.StaffRepository) service.StaffService {
	return &staffService{staffRepo: staffRepo}
}

func (s *staffService) GetStaffByID(ctx context.Context, id int64) (*model.Staff, error) {
	if id <= 0 {
		return nil, errors.New("invalid_staff_id")
	}
	return s.staffRepo.GetStaffByID(ctx, id)
}
