package service

import (
	"context"

	"ads-platform/internal/business/iam/action/model"
)

type ActionService interface {
	Create(ctx context.Context, code, name string) (*model.Action, error)
	Get(ctx context.Context, id int64) (*model.Action, error)
	List(ctx context.Context) ([]model.Action, error)
	Update(ctx context.Context, id int64, name string) (*model.Action, error)
	Delete(ctx context.Context, id int64) error
}
