package impl

import (
	"context"
	"errors"

	"ads-platform/internal/business/iam/action/errorcode"
	"ads-platform/internal/business/iam/action/model"
	"ads-platform/internal/business/iam/action/repository"
	"ads-platform/internal/business/iam/action/service"
	"ads-platform/internal/business/iam/internal/named"
	"ads-platform/internal/business/iam/internal/pgerr"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type actionService struct {
	repo repository.ActionRepository
}

func NewActionService(repo repository.ActionRepository) service.ActionService {
	return &actionService{repo: repo}
}

func (s *actionService) Create(ctx context.Context, code, name string) (*model.Action, error) {
	code, name, err := normalizeNamed(code, name)
	if err != nil {
		return nil, err
	}
	action := &model.Action{Code: code, Name: name}
	if err := s.repo.Create(ctx, action); err != nil {
		if pgerr.IsUniqueViolation(err) {
			return nil, exception.NewAppError(errorcode.ErrDuplicateCode.Code, errorcode.ErrDuplicateCode.HttpStatus)
		}
		return nil, err
	}
	return action, nil
}

func (s *actionService) Get(ctx context.Context, id int64) (*model.Action, error) {
	if id <= 0 {
		return nil, exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	action, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
		}
		return nil, err
	}
	return action, nil
}

func (s *actionService) List(ctx context.Context) ([]model.Action, error) {
	return s.repo.List(ctx)
}

func (s *actionService) Update(ctx context.Context, id int64, name string) (*model.Action, error) {
	action, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	name, ok := named.Name(name)
	if !ok {
		return nil, exception.NewAppError(errorcode.ErrInvalidName.Code, errorcode.ErrInvalidName.HttpStatus)
	}
	action.Name = name
	if err := s.repo.Update(ctx, action); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *actionService) Delete(ctx context.Context, id int64) error {
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

func normalizeNamed(code, name string) (string, string, error) {
	code, ok := named.Code(code)
	if !ok {
		return "", "", exception.NewAppError(errorcode.ErrInvalidCode.Code, errorcode.ErrInvalidCode.HttpStatus, code)
	}
	name, ok = named.Name(name)
	if !ok {
		return "", "", exception.NewAppError(errorcode.ErrInvalidName.Code, errorcode.ErrInvalidName.HttpStatus)
	}
	return code, name, nil
}
