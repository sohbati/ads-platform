package service

import (
	"context"

	"ads-platform/internal/business/iam/role/model"
)

type RoleService interface {
	Create(ctx context.Context, code, name string) (*model.Role, error)
	Get(ctx context.Context, id int64) (*model.Role, error)
	List(ctx context.Context) ([]model.Role, error)
	// Update renames the item and, when code is non-empty, changes its code too.
	Update(ctx context.Context, id int64, code, name string) (*model.Role, error)
	Delete(ctx context.Context, id int64) error
}
