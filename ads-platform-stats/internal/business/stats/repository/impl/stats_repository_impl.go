package impl

import (
	"context"
	"errors"
	"time"

	"ads-platform-stats/internal/business/stats/model"
	"ads-platform-stats/internal/business/stats/repository"

	"gorm.io/gorm"
)

type statsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) repository.StatsRepository {
	return &statsRepository{db: db}
}

func (r *statsRepository) AdOwnerID(ctx context.Context, adID int64) (int64, bool, error) {
	var userID int64
	err := r.db.WithContext(ctx).
		Table("ads").
		Select("user_id").
		Where("id = ? AND status <> ?", adID, "deleted").
		Take(&userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return userID, true, nil
}

func (r *statsRepository) Record(ctx context.Context, ev model.Event, day time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			`INSERT INTO ad_events (ad_id, event, viewer_id, session_user_id, occurred_at)
			 VALUES (?, ?, ?, ?, ?)`,
			ev.AdID, ev.Event, ev.ViewerID, ev.SessionUserID, ev.OccurredAt,
		).Error; err != nil {
			return err
		}

		views, unique, reveals, calls := 0, 0, 0, 0
		switch ev.Event {
		case model.EventView:
			views = 1
			res := tx.Exec(
				`INSERT INTO ad_view_uniques (ad_id, viewer_id, day) VALUES (?, ?, ?)
				 ON CONFLICT (ad_id, viewer_id, day) DO NOTHING`,
				ev.AdID, ev.ViewerID, day,
			)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected > 0 {
				unique = 1
			}
		case model.EventContactReveal:
			reveals = 1
		case model.EventCall:
			calls = 1
		}

		return tx.Exec(
			`INSERT INTO ad_stats_daily (ad_id, day, views, unique_viewers, contact_reveals, calls)
			 VALUES (?, ?, ?, ?, ?, ?)
			 ON CONFLICT (ad_id, day) DO UPDATE SET
			   views = ad_stats_daily.views + EXCLUDED.views,
			   unique_viewers = ad_stats_daily.unique_viewers + EXCLUDED.unique_viewers,
			   contact_reveals = ad_stats_daily.contact_reveals + EXCLUDED.contact_reveals,
			   calls = ad_stats_daily.calls + EXCLUDED.calls`,
			ev.AdID, day, views, unique, reveals, calls,
		).Error
	})
}
