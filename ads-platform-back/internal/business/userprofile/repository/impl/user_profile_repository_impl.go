package impl

import (
	"context"

	"ads-platform/internal/business/userprofile/model"
	"ads-platform/internal/business/userprofile/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) repository.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) GetByUserID(ctx context.Context, userID int64) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userProfileRepository) Upsert(ctx context.Context, profile *model.UserProfile) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"location_slugs", "updated_at"}),
	}).Create(profile).Error
}
