package repository

import (
	"context"

	"ads-platform/internal/business/ads/model"
)

type AdImageRepository interface {
	Create(ctx context.Context, image *model.AdImage) error
	GetByID(ctx context.Context, id int64) (*model.AdImage, error)
	Update(ctx context.Context, image *model.AdImage) error
	// NextObjectSeq is 1 + the number of images this customer has ever stored.
	NextObjectSeq(ctx context.Context, customerID int64) (int64, error)
}
