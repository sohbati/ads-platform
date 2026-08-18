package service

import (
	"context"

	"cache-service/internal/simplecache/category/model"
)

type CategoryCacheService interface {
	GetBySlugs(ctx context.Context, slugs []string) ([]model.Category, error)
	GetByIDs(ctx context.Context, ids []int) ([]model.Category, error)
}
