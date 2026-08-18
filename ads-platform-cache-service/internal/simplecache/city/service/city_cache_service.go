package service

import (
	"context"

	"cache-service/internal/simplecache/city/model"
)

type CityCacheService interface {
	GetByIDs(ctx context.Context, ids []int) ([]model.City, error)
	GetBySlugs(ctx context.Context, slugs []string) ([]model.City, error)
}
