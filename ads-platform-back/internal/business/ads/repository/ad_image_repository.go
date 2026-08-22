package repository

import (
	"context"

	"ads-platform/internal/business/ads/model"
)

type AdImageRepository interface {
	Create(ctx context.Context, image *model.AdImage) error
	GetByID(ctx context.Context, id int64) (*model.AdImage, error)
	Update(ctx context.Context, image *model.AdImage) error
}
