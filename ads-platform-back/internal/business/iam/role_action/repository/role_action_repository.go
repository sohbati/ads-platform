package repository

import (
	"context"

	actionmodel "ads-platform/internal/business/iam/action/model"
)

type RoleActionRepository interface {
	ListActions(ctx context.Context, roleID int64) ([]actionmodel.Action, error)
	Replace(ctx context.Context, roleID int64, actionIDs []int64) error
	Add(ctx context.Context, roleID, actionID int64) error
	Remove(ctx context.Context, roleID, actionID int64) error
}
