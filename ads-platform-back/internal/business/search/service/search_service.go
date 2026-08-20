package service

import (
	"context"

	"ads-platform/internal/business/search/model"
)

type SearchService interface {
	Search(ctx context.Context, placeSlug, categorySlug string, citiesCSV string, query string, priceMin, priceMax *int64, hasPhoto *bool, neighborhood, sort string, page, limit int) (*model.SearchResponse, error)
}
