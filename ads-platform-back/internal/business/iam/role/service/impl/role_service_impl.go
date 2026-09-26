package impl

import (
	"context"
	"errors"

	"ads-platform/internal/business/iam/internal/named"
	"ads-platform/internal/business/iam/internal/pgerr"
	"ads-platform/internal/business/iam/role/errorcode"
	"ads-platform/internal/business/iam/role/model"
	"ads-platform/internal/business/iam/role/repository"
	"ads-platform/internal/business/iam/role/service"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) service.RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) Create(ctx context.Context, code, name string) (*model.Role, error) {
	code, ok := named.Code(code)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidCode.Code, errorcode.ErrInvalidCode.HttpStatus, code)
	}
	name, ok = named.Name(name)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidName.Code, errorcode.ErrInvalidName.HttpStatus)
	}
	role := &model.Role{Code: code, Name: name}
	if err := s.repo.Create(ctx, role); err != nil {
		if pgerr.IsUniqueViolation(err) {
			return nil, exception.NewAppError(errorcode.ErrDuplicateCode.Code, errorcode.ErrDuplicateCode.HttpStatus)
		}
		return nil, err
	}
	return role, nil
}

func (s *roleService) Get(ctx context.Context, id int64) (*model.Role, error) {
	if id <= 0 {
		return nil, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
		}
		return nil, err
	}
	return role, nil
}

func (s *roleService) List(ctx context.Context) ([]model.Role, error) {
	return s.repo.List(ctx)
}

func (s *roleService) Update(ctx context.Context, id int64, name string) (*model.Role, error) {
	role, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	name, ok := named.Name(name)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidName.Code, errorcode.ErrInvalidName.HttpStatus)
	}
	role.Name = name
	if err := s.repo.Update(ctx, role); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *roleService) Delete(ctx context.Context, id int64) error {
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
