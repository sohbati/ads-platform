package impl

import (
	"context"
	"encoding/json"
	"time"

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

func (r *adRepository) GetByID(ctx context.Context, id int64) (*model.Ad, error) {
	var ad model.Ad
	if err := r.db.WithContext(ctx).First(&ad, id).Error; err != nil {
		return nil, err
	}
	return &ad, nil
}

func (r *adRepository) Update(ctx context.Context, ad *model.Ad) error {
	return r.db.WithContext(ctx).Model(ad).Select(
		"category_id", "city_id", "title", "description",
		"price_amount", "price_type", "currency",
		"attrs", "contact", "location", "media",
	).Updates(ad).Error
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

func (r *adRepository) ListStats(ctx context.Context, userID int64, from, to time.Time) ([]model.AdStatsItem, error) {
	var items []model.AdStatsItem
	err := r.db.WithContext(ctx).Raw(
		`SELECT d.ad_id,
		        COALESCE(SUM(d.views), 0) AS views,
		        COALESCE(SUM(d.unique_viewers), 0) AS unique_viewers,
		        COALESCE(SUM(d.contact_reveals), 0) AS contact_reveals,
		        COALESCE(SUM(d.calls), 0) AS calls
		 FROM ad_stats_daily d
		 INNER JOIN ads a ON a.id = d.ad_id
		 WHERE a.user_id = ?
		   AND a.status <> ?
		   AND d.day >= ?
		   AND d.day <= ?
		 GROUP BY d.ad_id`,
		userID, model.AdStatusDeleted, from, to,
	).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []model.AdStatsItem{}
	}
	return items, nil
}
