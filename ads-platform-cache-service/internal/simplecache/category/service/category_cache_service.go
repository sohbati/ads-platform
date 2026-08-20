package service

import (
	"context"

	"cache-service/internal/simplecache/category/model"
)

type CategoryCacheService interface {
	GetBySlugs(ctx context.Context, slugs []string, includeDescendants bool) ([]model.Category, error)
	GetByIDs(ctx context.Context, ids []int, includeDescendants bool) ([]model.Category, error)
}
