package service

import (
	"context"

	actionmodel "ads-platform/internal/business/iam/action/model"
)

type RoleActionService interface {
	List(ctx context.Context, roleID int64) ([]actionmodel.Action, error)
	Set(ctx context.Context, roleID int64, actionIDs []int64) ([]actionmodel.Action, error)
	Add(ctx context.Context, roleID, actionID int64) ([]actionmodel.Action, error)
	Remove(ctx context.Context, roleID, actionID int64) error
}
