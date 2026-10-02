package service

import (
	"context"

	"ads-platform/internal/business/customerprofile/model"
)

type CustomerProfileService interface {
	Get(ctx context.Context, customerID int64) (*model.CustomerProfile, error)
	Put(ctx context.Context, customerID int64, locationSlugs []string) (*model.CustomerProfile, error)
}
