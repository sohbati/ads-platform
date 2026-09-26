package repository

import (
	"context"

	"ads-platform/internal/business/iam/user_group/model"
)

type UserGroupRepository interface {
	Create(ctx context.Context, group *model.UserGroup) error
	GetByID(ctx context.Context, id int64) (*model.UserGroup, error)
	List(ctx context.Context) ([]model.UserGroup, error)
	Update(ctx context.Context, group *model.UserGroup) error
	Delete(ctx context.Context, id int64) error
	ListUserIDs(ctx context.Context, groupID int64) ([]int64, error)
	AddUser(ctx context.Context, groupID, userID int64) error
	RemoveUser(ctx context.Context, groupID, userID int64) error
	ListByUser(ctx context.Context, userID int64) ([]model.UserGroup, error)
}
