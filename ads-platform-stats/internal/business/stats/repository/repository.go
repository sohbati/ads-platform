package repository

import (
	"context"
	"time"

	"ads-platform-stats/internal/business/stats/model"
)

type StatsRepository interface {
	AdOwnerID(ctx context.Context, adID int64) (int64, bool, error)
	Record(ctx context.Context, ev model.Event, day time.Time) error
}
