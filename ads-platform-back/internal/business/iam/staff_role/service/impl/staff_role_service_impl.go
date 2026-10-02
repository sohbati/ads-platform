package impl

import (
	"context"
	"errors"
	"strings"

	actionmodel "ads-platform/internal/business/iam/action/model"
	"ads-platform/internal/business/iam/internal/named"
	"ads-platform/internal/business/iam/internal/pgerr"
	rolemodel "ads-platform/internal/business/iam/role/model"
	staffmodel "ads-platform/internal/business/iam/staff/model"
	"ads-platform/internal/business/iam/staff_role/errorcode"
	"ads-platform/internal/business/iam/staff_role/repository"
	"ads-platform/internal/business/iam/staff_role/service"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type staffLookup interface {
	GetStaffByID(ctx context.Context, id int64) (*staffmodel.Staff, error)
}

type roleLookup interface {
	GetByID(ctx context.Context, id int64) (*rolemodel.Role, error)
}

type staffRoleService struct {
	repo  repository.StaffRoleRepository
	staff staffLookup
	roles roleLookup
}

func NewStaffRoleService(repo repository.StaffRoleRepository, staff staffLookup, roles roleLookup) service.StaffRoleService {
	return &staffRoleService{repo: repo, staff: staff, roles: roles}
}

func (s *staffRoleService) ListRoles(ctx context.Context, staffID int64) ([]rolemodel.Role, error) {
	if err := s.ensureStaff(ctx, staffID); err != nil {
		return nil, err
	}
	return s.repo.ListRoles(ctx, staffID)
}

func (s *staffRoleService) Add(ctx context.Context, staffID, roleID int64) ([]rolemodel.Role, error) {
	if err := s.ensureStaff(ctx, staffID); err != nil {
		return nil, err
	}
	if err := s.ensureRole(ctx, roleID); err != nil {
		return nil, err
	}
	if err := s.repo.Add(ctx, staffID, roleID); err != nil && !pgerr.IsUniqueViolation(err) {
		return nil, err
	}
	return s.repo.ListRoles(ctx, staffID)
}

func (s *staffRoleService) Remove(ctx context.Context, staffID, roleID int64) error {
	if err := s.ensureStaff(ctx, staffID); err != nil {
		return err
	}
	if err := s.repo.Remove(ctx, staffID, roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
		}
		return err
	}
	return nil
}

func (s *staffRoleService) ListActions(ctx context.Context, staffID int64) ([]actionmodel.Action, error) {
	if err := s.ensureStaff(ctx, staffID); err != nil {
		return nil, err
	}
	return s.repo.ListActions(ctx, staffID)
}

func (s *staffRoleService) HasAction(ctx context.Context, staffID int64, code string) (bool, error) {
	if err := s.ensureStaff(ctx, staffID); err != nil {
		return false, err
	}
	code = strings.ToLower(strings.TrimSpace(code))
	if !named.CodePattern.MatchString(code) {
		return false, exception.NewAppError(errorcode.ErrInvalidCode.Code, errorcode.ErrInvalidCode.HttpStatus, code)
	}
	return s.repo.HasAction(ctx, staffID, code)
}

func (s *staffRoleService) RequireAction(ctx context.Context, staffID int64, code string) error {
	ok, err := s.HasAction(ctx, staffID, code)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewAppError(errorcode.ErrActionDenied.Code, errorcode.ErrActionDenied.HttpStatus, code)
	}
	return nil
}

func (s *staffRoleService) ensureStaff(ctx context.Context, staffID int64) error {
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

func (s *staffRoleService) ensureRole(ctx context.Context, roleID int64) error {
	if roleID <= 0 {
		return exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	_, err := s.roles.GetByID(ctx, roleID)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.NewAppError(errorcode.ErrRoleNotFound.Code, errorcode.ErrRoleNotFound.HttpStatus)
	}
	return err
}
