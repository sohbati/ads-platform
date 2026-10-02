package service

import (
	"context"

	actionmodel "ads-platform/internal/business/iam/action/model"
	rolemodel "ads-platform/internal/business/iam/role/model"
)

type StaffRoleService interface {
	ListRoles(ctx context.Context, staffID int64) ([]rolemodel.Role, error)
	Add(ctx context.Context, staffID, roleID int64) ([]rolemodel.Role, error)
	Remove(ctx context.Context, staffID, roleID int64) error
	ListActions(ctx context.Context, staffID int64) ([]actionmodel.Action, error)
	HasAction(ctx context.Context, staffID int64, code string) (bool, error)
	RequireAction(ctx context.Context, staffID int64, code string) error
}
