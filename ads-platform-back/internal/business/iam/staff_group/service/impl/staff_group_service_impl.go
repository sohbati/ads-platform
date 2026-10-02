package impl

import (
	"context"
	"errors"

	"ads-platform/internal/business/iam/internal/named"
	"ads-platform/internal/business/iam/internal/pgerr"
	staffmodel "ads-platform/internal/business/iam/staff/model"
	"ads-platform/internal/business/iam/staff_group/errorcode"
	"ads-platform/internal/business/iam/staff_group/model"
	"ads-platform/internal/business/iam/staff_group/repository"
	"ads-platform/internal/business/iam/staff_group/service"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type staffLookup interface {
	GetStaffByID(ctx context.Context, id int64) (*staffmodel.Staff, error)
}

type staffGroupService struct {
	repo  repository.StaffGroupRepository
	staff staffLookup
}

func NewStaffGroupService(repo repository.StaffGroupRepository, staff staffLookup) service.StaffGroupService {
	return &staffGroupService{repo: repo, staff: staff}
}

func (s *staffGroupService) Create(ctx context.Context, code, name string) (*model.StaffGroup, error) {
	code, ok := named.Code(code)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidCode.Code, errorcode.ErrInvalidCode.HttpStatus, code)
	}
	name, ok = named.Name(name)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidName.Code, errorcode.ErrInvalidName.HttpStatus)
	}
	group := &model.StaffGroup{Code: code, Name: name}
	if err := s.repo.Create(ctx, group); err != nil {
		if pgerr.IsUniqueViolation(err) {
			return nil, exception.NewAppError(errorcode.ErrDuplicateCode.Code, errorcode.ErrDuplicateCode.HttpStatus)
		}
		return nil, err
	}
	return group, nil
}

func (s *staffGroupService) Get(ctx context.Context, id int64) (*model.StaffGroup, error) {
	if id <= 0 {
		return nil, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	group, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
		}
		return nil, err
	}
	return group, nil
}

func (s *staffGroupService) List(ctx context.Context) ([]model.StaffGroup, error) {
	return s.repo.List(ctx)
}

func (s *staffGroupService) Update(ctx context.Context, id int64, name string) (*model.StaffGroup, error) {
	group, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	name, ok := named.Name(name)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidName.Code, errorcode.ErrInvalidName.HttpStatus)
	}
	group.Name = name
	if err := s.repo.Update(ctx, group); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *staffGroupService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
		}
		return err
	}
	return nil
}

func (s *staffGroupService) ListStaff(ctx context.Context, groupID int64) ([]int64, error) {
	if _, err := s.Get(ctx, groupID); err != nil {
		return nil, err
	}
	return s.repo.ListStaffIDs(ctx, groupID)
}

func (s *staffGroupService) AddStaff(ctx context.Context, groupID, staffID int64) error {
	if _, err := s.Get(ctx, groupID); err != nil {
		return err
	}
	if err := s.ensureStaff(ctx, staffID); err != nil {
		return err
	}
	if err := s.repo.AddStaff(ctx, groupID, staffID); err != nil && !pgerr.IsUniqueViolation(err) {
		return err
	}
	return nil
}

func (s *staffGroupService) RemoveStaff(ctx context.Context, groupID, staffID int64) error {
	if _, err := s.Get(ctx, groupID); err != nil {
		return err
	}
	if err := s.repo.RemoveStaff(ctx, groupID, staffID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewAppError(errorcode.ErrMembershipNotFound.Code, errorcode.ErrMembershipNotFound.HttpStatus)
		}
		return err
	}
	return nil
}

func (s *staffGroupService) ListByStaff(ctx context.Context, staffID int64) ([]model.StaffGroup, error) {
	if err := s.ensureStaff(ctx, staffID); err != nil {
		return nil, err
	}
	return s.repo.ListByStaff(ctx, staffID)
}

func (s *staffGroupService) ensureStaff(ctx context.Context, staffID int64) error {
	if staffID <= 0 {
		return exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	_, err := s.staff.GetStaffByID(ctx, staffID)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.NewAppError(errorcode.ErrStaffNotFound.Code, errorcode.ErrStaffNotFound.HttpStatus)
	}
	return err
}
