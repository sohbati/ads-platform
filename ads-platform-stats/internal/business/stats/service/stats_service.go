package service

import (
	"context"
	"time"

	"ads-platform-stats/internal/business/stats/model"
	"ads-platform-stats/internal/business/stats/repository"
)

type StatsService interface {
	Apply(ctx context.Context, ev model.Event) error
}

type statsService struct {
	repo repository.StatsRepository
	loc  *time.Location
}

func NewStatsService(repo repository.StatsRepository) StatsService {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		loc = time.FixedZone("Tehran", 3*3600+30*60)
	}
	return &statsService{repo: repo, loc: loc}
}

func (s *statsService) Apply(ctx context.Context, ev model.Event) error {
	if ev.AdID <= 0 || ev.ViewerID == "" || !model.ValidEvent(ev.Event) {
		return nil
	}
	if ev.OccurredAt.IsZero() {
		ev.OccurredAt = time.Now().UTC()
	}

	ownerID, ok, err := s.repo.AdOwnerID(ctx, ev.AdID)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	if ev.SessionUserID != nil && *ev.SessionUserID == ownerID {
		return nil
	}

	day := ev.OccurredAt.In(s.loc)
	day = time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return s.repo.Record(ctx, ev, day)
}
