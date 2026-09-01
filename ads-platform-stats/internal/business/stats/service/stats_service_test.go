package service_test

import (
	"context"
	"testing"
	"time"

	"ads-platform-stats/internal/business/stats/model"
	"ads-platform-stats/internal/business/stats/service"
)

type fakeRepo struct {
	ownerID   int64
	found     bool
	ownerErr  error
	recorded  []model.Event
	recordErr error
}

func (f *fakeRepo) AdOwnerID(context.Context, int64) (int64, bool, error) {
	return f.ownerID, f.found, f.ownerErr
}

func (f *fakeRepo) Record(_ context.Context, ev model.Event, _ time.Time) error {
	f.recorded = append(f.recorded, ev)
	return f.recordErr
}

func TestApplyDropsOwnerViews(t *testing.T) {
	repo := &fakeRepo{ownerID: 7, found: true}
	svc := service.NewStatsService(repo)
	uid := int64(7)
	err := svc.Apply(context.Background(), model.Event{
		AdID: 55, Event: model.EventView, ViewerID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		SessionUserID: &uid, OccurredAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.recorded) != 0 {
		t.Fatalf("owner event should be dropped, got %+v", repo.recorded)
	}
}

func TestApplyRecordsGuestView(t *testing.T) {
	repo := &fakeRepo{ownerID: 7, found: true}
	svc := service.NewStatsService(repo)
	err := svc.Apply(context.Background(), model.Event{
		AdID: 55, Event: model.EventView, ViewerID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		OccurredAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(repo.recorded) != 1 || repo.recorded[0].Event != model.EventView {
		t.Fatalf("recorded=%+v", repo.recorded)
	}
}

func TestApplyDropsUnknownAd(t *testing.T) {
	repo := &fakeRepo{found: false}
	svc := service.NewStatsService(repo)
	_ = svc.Apply(context.Background(), model.Event{
		AdID: 99, Event: model.EventView, ViewerID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
	})
	if len(repo.recorded) != 0 {
		t.Fatal("unknown ad should be dropped")
	}
}

type tallyRepo struct {
	ownerID int64
	found   bool
	views   int
	unique  int
	seen    map[string]struct{}
}

func (t *tallyRepo) AdOwnerID(context.Context, int64) (int64, bool, error) {
	return t.ownerID, t.found, nil
}

func (t *tallyRepo) Record(_ context.Context, ev model.Event, day time.Time) error {
	if t.seen == nil {
		t.seen = map[string]struct{}{}
	}
	if ev.Event != model.EventView {
		return nil
	}
	t.views++
	key := ev.ViewerID + "|" + day.Format("2006-01-02")
	if _, ok := t.seen[key]; !ok {
		t.seen[key] = struct{}{}
		t.unique++
	}
	return nil
}

func TestApplySameViewerSameDayCountsUniqueOnce(t *testing.T) {
	repo := &tallyRepo{ownerID: 7, found: true}
	svc := service.NewStatsService(repo)
	when := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	ev := model.Event{
		AdID: 55, Event: model.EventView, ViewerID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		OccurredAt: when,
	}
	if err := svc.Apply(context.Background(), ev); err != nil {
		t.Fatal(err)
	}
	ev.OccurredAt = when.Add(2 * time.Hour)
	if err := svc.Apply(context.Background(), ev); err != nil {
		t.Fatal(err)
	}
	if repo.views != 2 || repo.unique != 1 {
		t.Fatalf("views=%d unique=%d, want 2 and 1", repo.views, repo.unique)
	}
}
