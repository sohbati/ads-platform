package service

import (
	"context"

	"ads-platform/internal/business/iam/role/model"
)

type RoleService interface {
	Create(ctx context.Context, code, name string) (*model.Role, error)
	Get(ctx context.Context, id int64) (*model.Role, error)
	List(ctx context.Context) ([]model.Role, error)
	Update(ctx context.Context, id int64, name string) (*model.Role, error)
	Delete(ctx context.Context, id int64) error
}
