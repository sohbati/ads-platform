package service

import (
	"context"

	"ads-platform/internal/business/userprofile/model"
)

type UserProfileService interface {
	Get(ctx context.Context, userID int64) (*model.UserProfile, error)
	Put(ctx context.Context, userID int64, locationSlugs []string) (*model.UserProfile, error)
}
