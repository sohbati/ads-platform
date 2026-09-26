package repository

import (
	"context"

	"ads-platform/internal/business/iam/action/model"
)

type ActionRepository interface {
	Create(ctx context.Context, action *model.Action) error
	GetByID(ctx context.Context, id int64) (*model.Action, error)
	List(ctx context.Context) ([]model.Action, error)
	Update(ctx context.Context, action *model.Action) error
	Delete(ctx context.Context, id int64) error
}
