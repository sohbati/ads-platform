package impl

import (
	"context"
	"errors"

	"ads-platform/internal/business/iam/internal/named"
	"ads-platform/internal/business/iam/internal/pgerr"
	usermodel "ads-platform/internal/business/iam/user/model"
	"ads-platform/internal/business/iam/user_group/errorcode"
	"ads-platform/internal/business/iam/user_group/model"
	"ads-platform/internal/business/iam/user_group/repository"
	"ads-platform/internal/business/iam/user_group/service"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type userLookup interface {
	GetUserByID(ctx context.Context, id int64) (*usermodel.User, error)
}

type userGroupService struct {
	repo  repository.UserGroupRepository
	users userLookup
}

func NewUserGroupService(repo repository.UserGroupRepository, users userLookup) service.UserGroupService {
	return &userGroupService{repo: repo, users: users}
}

func (s *userGroupService) Create(ctx context.Context, code, name string) (*model.UserGroup, error) {
	code, ok := named.Code(code)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidCode.Code, errorcode.ErrInvalidCode.HttpStatus, code)
	}
	name, ok = named.Name(name)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidName.Code, errorcode.ErrInvalidName.HttpStatus)
	}
	group := &model.UserGroup{Code: code, Name: name}
	if err := s.repo.Create(ctx, group); err != nil {
		if pgerr.IsUniqueViolation(err) {
			return nil, exception.NewAppError(errorcode.ErrDuplicateCode.Code, errorcode.ErrDuplicateCode.HttpStatus)
		}
		return nil, err
	}
	return group, nil
}

func (s *userGroupService) Get(ctx context.Context, id int64) (*model.UserGroup, error) {
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

func (s *userGroupService) List(ctx context.Context) ([]model.UserGroup, error) {
	return s.repo.List(ctx)
}

func (s *userGroupService) Update(ctx context.Context, id int64, name string) (*model.UserGroup, error) {
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

func (s *userGroupService) Delete(ctx context.Context, id int64) error {
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

func (s *userGroupService) ListUsers(ctx context.Context, groupID int64) ([]int64, error) {
	if _, err := s.Get(ctx, groupID); err != nil {
		return nil, err
	}
	return s.repo.ListUserIDs(ctx, groupID)
}

func (s *userGroupService) AddUser(ctx context.Context, groupID, userID int64) error {
	if _, err := s.Get(ctx, groupID); err != nil {
		return err
	}
	if err := s.ensureUser(ctx, userID); err != nil {
		return err
	}
	if err := s.repo.AddUser(ctx, groupID, userID); err != nil && !pgerr.IsUniqueViolation(err) {
		return err
	}
	return nil
}

func (s *userGroupService) RemoveUser(ctx context.Context, groupID, userID int64) error {
	if _, err := s.Get(ctx, groupID); err != nil {
		return err
	}
	if err := s.repo.RemoveUser(ctx, groupID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewAppError(errorcode.ErrMembershipNotFound.Code, errorcode.ErrMembershipNotFound.HttpStatus)
		}
		return err
	}
	return nil
}

func (s *userGroupService) ListByUser(ctx context.Context, userID int64) ([]model.UserGroup, error) {
	if err := s.ensureUser(ctx, userID); err != nil {
		return nil, err
	}
	return s.repo.ListByUser(ctx, userID)
}

func (s *userGroupService) ensureUser(ctx context.Context, userID int64) error {
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
