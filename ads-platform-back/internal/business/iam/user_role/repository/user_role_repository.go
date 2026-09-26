package repository

import (
	"context"

	actionmodel "ads-platform/internal/business/iam/action/model"
	rolemodel "ads-platform/internal/business/iam/role/model"
)

type UserRoleRepository interface {
	ListRoles(ctx context.Context, userID int64) ([]rolemodel.Role, error)
	Add(ctx context.Context, userID, roleID int64) error
	Remove(ctx context.Context, userID, roleID int64) error
	ListActions(ctx context.Context, userID int64) ([]actionmodel.Action, error)
	HasAction(ctx context.Context, userID int64, code string) (bool, error)
}
