package repository

import (
	"context"
	"encoding/json"
	"time"

	"ads-platform/internal/business/ads/model"
)

type AdRepository interface {
	Create(ctx context.Context, ad *model.Ad) error
	GetByID(ctx context.Context, id int64) (*model.Ad, error)
	Update(ctx context.Context, ad *model.Ad) error
	UpdateMedia(ctx context.Context, id int64, media json.RawMessage) error
	ListByUserID(ctx context.Context, userID int64, limit int) ([]model.Ad, error)
	ListStats(ctx context.Context, userID int64, from, to time.Time) ([]model.AdStatsItem, error)
}
