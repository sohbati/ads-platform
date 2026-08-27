package repository

import (
	"context"
	"encoding/json"

	"ads-platform/internal/business/ads/model"
)

type AdRepository interface {
	Create(ctx context.Context, ad *model.Ad) error
	GetByID(ctx context.Context, id int64) (*model.Ad, error)
	Update(ctx context.Context, ad *model.Ad) error
	UpdateMedia(ctx context.Context, id int64, media json.RawMessage) error
	ListByUserID(ctx context.Context, userID int64, limit int) ([]model.Ad, error)
}
