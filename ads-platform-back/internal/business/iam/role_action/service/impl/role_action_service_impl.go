package impl

import (
	"context"
	"errors"

	actionmodel "ads-platform/internal/business/iam/action/model"
	actionrepo "ads-platform/internal/business/iam/action/repository"
	"ads-platform/internal/business/iam/internal/pgerr"
	rolemodel "ads-platform/internal/business/iam/role/model"
	rolerepo "ads-platform/internal/business/iam/role/repository"
	"ads-platform/internal/business/iam/role_action/errorcode"
	"ads-platform/internal/business/iam/role_action/repository"
	"ads-platform/internal/business/iam/role_action/service"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type roleLookup interface {
	GetByID(ctx context.Context, id int64) (*rolemodel.Role, error)
}

type actionLookup interface {
	GetByID(ctx context.Context, id int64) (*actionmodel.Action, error)
}

type roleActionService struct {
	repo    repository.RoleActionRepository
	roles   roleLookup
	actions actionLookup
}

func NewRoleActionService(repo repository.RoleActionRepository, roles rolerepo.RoleRepository, actions actionrepo.ActionRepository) service.RoleActionService {
	return &roleActionService{repo: repo, roles: roles, actions: actions}
}

func (s *roleActionService) List(ctx context.Context, roleID int64) ([]actionmodel.Action, error) {
	if err := s.ensureRole(ctx, roleID); err != nil {
		return nil, err
	}
	return s.repo.ListActions(ctx, roleID)
}

func (s *roleActionService) Set(ctx context.Context, roleID int64, actionIDs []int64) ([]actionmodel.Action, error) {
	if err := s.ensureRole(ctx, roleID); err != nil {
		return nil, err
	}
	ids, err := s.uniqueExistingActionIDs(ctx, actionIDs)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Replace(ctx, roleID, ids); err != nil {
		return nil, err
	}
	return s.repo.ListActions(ctx, roleID)
}

func (s *roleActionService) Add(ctx context.Context, roleID, actionID int64) ([]actionmodel.Action, error) {
	if err := s.ensureRole(ctx, roleID); err != nil {
		return nil, err
	}
	if err := s.ensureAction(ctx, actionID); err != nil {
		return nil, err
	}
	if err := s.repo.Add(ctx, roleID, actionID); err != nil && !pgerr.IsUniqueViolation(err) {
		return nil, err
	}
	return s.repo.ListActions(ctx, roleID)
}

func (s *roleActionService) Remove(ctx context.Context, roleID, actionID int64) error {
	if err := s.ensureRole(ctx, roleID); err != nil {
		return err
	}
	if err := s.repo.Remove(ctx, roleID, actionID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewAppError(errorcode.ErrNotFound.Code, errorcode.ErrNotFound.HttpStatus)
		}
		return err
	}
	return nil
}

func (s *roleActionService) ensureRole(ctx context.Context, roleID int64) error {
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

func (s *roleActionService) ensureAction(ctx context.Context, actionID int64) error {
	if actionID <= 0 {
		return exception.NewAppError(errorcode.ErrInvalidID.Code, errorcode.ErrInvalidID.HttpStatus)
	}
	_, err := s.actions.GetByID(ctx, actionID)
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exception.NewAppError(errorcode.ErrActionNotFound.Code, errorcode.ErrActionNotFound.HttpStatus)
	}
	return err
}

func (s *roleActionService) uniqueExistingActionIDs(ctx context.Context, actionIDs []int64) ([]int64, error) {
	out := make([]int64, 0, len(actionIDs))
	seen := make(map[int64]struct{}, len(actionIDs))
	for _, id := range actionIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		if err := s.ensureAction(ctx, id); err != nil {
			return nil, err
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out, nil
}
