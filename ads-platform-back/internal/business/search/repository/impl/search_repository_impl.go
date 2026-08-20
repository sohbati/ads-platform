package impl

import (
	"context"
	"strings"

	adsmodel "ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/search/model"
	"ads-platform/internal/business/search/repository"

	"gorm.io/gorm"
)

type searchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) repository.SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) SearchAds(ctx context.Context, filter model.SearchFilter) ([]adsmodel.Ad, int64, error) {
	q := r.db.WithContext(ctx).Model(&adsmodel.Ad{}).Where("status = ?", "active")

	if len(filter.CategoryIDs) > 0 {
		q = q.Where("category_id IN ?", filter.CategoryIDs)
	}
	if len(filter.CityIDs) > 0 {
		q = q.Where("city_id IN ?", filter.CityIDs)
	}
	if filter.Query != "" {
		term := "%" + filter.Query + "%"
		q = q.Where("(title ILIKE ? OR description ILIKE ?)", term, term)
	}
	if filter.PriceMin != nil {
		q = q.Where("price_amount >= ?", *filter.PriceMin)
	}
	if filter.PriceMax != nil {
		q = q.Where("price_amount <= ?", *filter.PriceMax)
	}
	if filter.HasPhoto != nil && *filter.HasPhoto {
		q = q.Where("jsonb_array_length(media) > 0")
	}
	if filter.Neighborhood != "" {
		q = q.Where("location->>'neighborhood' ILIKE ?", "%"+filter.Neighborhood+"%")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch strings.ToLower(filter.Sort) {
	case "cheapest":
		q = q.Order("price_amount ASC NULLS LAST")
	case "expensive":
		q = q.Order("price_amount DESC NULLS LAST")
	default:
		q = q.Order("published_at DESC NULLS LAST").Order("id DESC")
	}

	offset := (filter.Page - 1) * filter.Limit
	var ads []adsmodel.Ad
	if err := q.Limit(filter.Limit).Offset(offset).Find(&ads).Error; err != nil {
		return nil, 0, err
	}
	return ads, total, nil
}
