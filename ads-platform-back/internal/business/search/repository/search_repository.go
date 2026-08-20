package repository

import (
	"context"

	adsmodel "ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/search/model"
)

type SearchRepository interface {
	SearchAds(ctx context.Context, filter model.SearchFilter) ([]adsmodel.Ad, int64, error)
}
