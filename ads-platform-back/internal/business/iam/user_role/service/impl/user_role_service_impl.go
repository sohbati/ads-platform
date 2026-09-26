package impl

import (
	"context"
	"errors"
	"strings"

	actionmodel "ads-platform/internal/business/iam/action/model"
	"ads-platform/internal/business/iam/internal/named"
	"ads-platform/internal/business/iam/internal/pgerr"
	rolemodel "ads-platform/internal/business/iam/role/model"
	usermodel "ads-platform/internal/business/iam/user/model"
	"ads-platform/internal/business/iam/user_role/errorcode"
	"ads-platform/internal/business/iam/user_role/repository"
	"ads-platform/internal/business/iam/user_role/service"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type userLookup interface {
	GetUserByID(ctx context.Context, id int64) (*usermodel.User, error)
}

type roleLookup interface {
	GetByID(ctx context.Context, id int64) (*rolemodel.Role, error)
}

type userRoleService struct {
	repo  repository.UserRoleRepository
	users userLookup
	roles roleLookup
}

func NewUserRoleService(repo repository.UserRoleRepository, users userLookup, roles roleLookup) service.UserRoleService {
	return &userRoleService{repo: repo, users: users, roles: roles}
}

func (s *userRoleService) ListRoles(ctx context.Context, userID int64) ([]rolemodel.Role, error) {
	if err := s.ensureUser(ctx, userID); err != nil {
		return nil, err
	}
	return s.repo.ListRoles(ctx, userID)
}

func (s *userRoleService) Add(ctx context.Context, userID, roleID int64) ([]rolemodel.Role, error) {
	if err := s.ensureUser(ctx, userID); err != nil {
		return nil, err
	}
	if err := s.ensureRole(ctx, roleID); err != nil {
		return nil, err
	}
	if err := s.repo.Add(ctx, userID, roleID); err != nil && !pgerr.IsUniqueViolation(err) {
		return nil, err
	}
	return s.repo.ListRoles(ctx, userID)
}

func (s *userRoleService) Remove(ctx context.Context, userID, roleID int64) error {
	if err := s.ensureUser(ctx, userID); err != nil {
		return err
	}
	if err := s.repo.Remove(ctx, userID, roleID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
		}
		return err
	}
	return nil
}

func (s *userRoleService) ListActions(ctx context.Context, userID int64) ([]actionmodel.Action, error) {
	if err := s.ensureUser(ctx, userID); err != nil {
		return nil, err
	}
	return s.repo.ListActions(ctx, userID)
}

func (s *userRoleService) HasAction(ctx context.Context, userID int64, code string) (bool, error) {
	if err := s.ensureUser(ctx, userID); err != nil {
		return false, err
	}
	code = strings.ToLower(strings.TrimSpace(code))
	if !named.CodePattern.MatchString(code) {
		return false, exception.NewAppError(errorcode.ErrInvalidCode.Code, errorcode.ErrInvalidCode.HttpStatus, code)
	}
	return s.repo.HasAction(ctx, userID, code)
}

func (s *userRoleService) RequireAction(ctx context.Context, userID int64, code string) error {
	ok, err := s.HasAction(ctx, userID, code)
	if err != nil {
		return err
	}
	if !ok {
		return exception.NewAppError(errorcode.ErrActionDenied.Code, errorcode.ErrActionDenied.HttpStatus, code)
	}
	return nil
}

func (s *userRoleService) ensureUser(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	_, err := s.users.GetUserByID(ctx, userID)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.NewAppError(errorcode.ErrUserNotFound.Code, errorcode.ErrUserNotFound.HttpStatus)
	}
	return err
}

func (s *userRoleService) ensureRole(ctx context.Context, roleID int64) error {
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
