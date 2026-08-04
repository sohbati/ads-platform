package service

import (
	"ads-platform/internal/business/user/model"
	"context"
)

// UserService defines the interface for user business logic operations
type UserService interface {
	// Basic CRUD operations
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int64) error

	// Search and filtering operations
	GetUserByMobile(ctx context.Context, mobile string) (*model.User, error)

	ValidateUserData(ctx context.Context, user *model.User) error
}
