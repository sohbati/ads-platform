package impl

import (
	"context"
	"encoding/json"

	"ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/ads/repository"

	"gorm.io/gorm"
)

type adRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) repository.AdRepository {
	return &adRepository{db: db}
}

func (r *adRepository) Create(ctx context.Context, ad *model.Ad) error {
	return r.db.WithContext(ctx).Create(ad).Error
}

func (r *adRepository) UpdateMedia(ctx context.Context, id int64, media json.RawMessage) error {
	return r.db.WithContext(ctx).Model(&model.Ad{}).Where("id = ?", id).Update("media", media).Error
}

func (r *adRepository) ListByUserID(ctx context.Context, userID int64, limit int) ([]model.Ad, error) {
	if limit < 1 {
		limit = 200
	}
	var ads []model.Ad
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("status <> ?", model.AdStatusDeleted).
		Order("published_at DESC NULLS LAST").
		Order("created_at DESC").
		Order("id DESC").
		Limit(limit).
		Find(&ads).Error
	return ads, err
}
