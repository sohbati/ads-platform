package impl

import (
	"context"

	"ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/ads/repository"

	"gorm.io/gorm"
)

type adImageRepository struct {
	db *gorm.DB
}

func NewAdImageRepository(db *gorm.DB) repository.AdImageRepository {
	return &adImageRepository{db: db}
}

func (r *adImageRepository) Create(ctx context.Context, image *model.AdImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *adImageRepository) GetByID(ctx context.Context, id int64) (*model.AdImage, error) {
	var image model.AdImage
	if err := r.db.WithContext(ctx).First(&image, id).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *adImageRepository) Update(ctx context.Context, image *model.AdImage) error {
	return r.db.WithContext(ctx).Save(image).Error
}

func (r *adImageRepository) NextObjectSeq(ctx context.Context, userID int64) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.AdImage{}).Where("user_id = ?", userID).Count(&n).Error
	if err != nil {
		return 0, err
	}
	return n + 1, nil
}
