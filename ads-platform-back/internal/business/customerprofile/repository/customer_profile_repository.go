package repository

import (
	"context"

	"ads-platform/internal/business/customerprofile/model"
)

type CustomerProfileRepository interface {
	GetByCustomerID(ctx context.Context, customerID int64) (*model.CustomerProfile, error)
	Upsert(ctx context.Context, profile *model.CustomerProfile) error
}
