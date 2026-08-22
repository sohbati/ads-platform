package repository

import (
	"context"
	"encoding/json"

	"ads-platform/internal/business/ads/model"
)

type AdRepository interface {
	Create(ctx context.Context, ad *model.Ad) error
	UpdateMedia(ctx context.Context, id int64, media json.RawMessage) error
}
