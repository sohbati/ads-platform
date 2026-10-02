package impl

import (
	"context"

	"ads-platform/internal/business/customerprofile/model"
	"ads-platform/internal/business/customerprofile/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type customerProfileRepository struct {
	db *gorm.DB
}

func NewCustomerProfileRepository(db *gorm.DB) repository.CustomerProfileRepository {
	return &customerProfileRepository{db: db}
}

func (r *customerProfileRepository) GetByCustomerID(ctx context.Context, customerID int64) (*model.CustomerProfile, error) {
	var profile model.CustomerProfile
	err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *customerProfileRepository) Upsert(ctx context.Context, profile *model.CustomerProfile) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "customer_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"location_slugs", "updated_at"}),
	}).Create(profile).Error
}
