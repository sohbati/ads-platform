package repository

import (
	"context"

	"ads-platform/internal/business/userprofile/model"
)

type UserProfileRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*model.UserProfile, error)
	Upsert(ctx context.Context, profile *model.UserProfile) error
}
