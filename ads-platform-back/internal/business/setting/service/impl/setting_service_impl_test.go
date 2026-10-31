package impl

import (
	"context"
	"testing"

	"ads-platform/internal/business/setting/model"
	"ads-platform/internal/core/exception"
)

type fakeRepo map[string]model.AppSetting

func (f fakeRepo) List(context.Context) ([]model.AppSetting, error) {
	out := []model.AppSetting{}
	for _, item := range f {
		out = append(out, item)
	}
	return out, nil
}

func (f fakeRepo) Get(_ context.Context, key string) (*model.AppSetting, error) {
	if item, ok := f[key]; ok {
		return &item, nil
	}
	return nil, nil
}

func (f fakeRepo) Upsert(_ context.Context, setting *model.AppSetting) error {
	f[setting.Key] = *setting
	return nil
}

func code(t *testing.T, err error) string {
	t.Helper()
	app, ok := exception.AsAppError(err)
	if !ok {
		t.Fatalf("want AppError, got %v", err)
	}
	return app.ErrorCode
}

func TestPendingDaysDefaultsToSeven(t *testing.T) {
	svc := NewSettingService(fakeRepo{})
	days, err := svc.Int(context.Background(), model.AdsPendingMaxDays)
	if err != nil || days != 7 {
		t.Fatalf("got %d, %v", days, err)
	}
	items, err := svc.List(context.Background())
	if err != nil || len(items) != 1 || items[0].Value != "7" || items[0].UpdatedAt != nil {
		t.Fatalf("got %+v, %v", items, err)
	}
}

func TestUpdatePendingDays(t *testing.T) {
	repo := fakeRepo{}
	svc := NewSettingService(repo)
	if _, err := svc.Update(context.Background(), model.AdsPendingMaxDays, " 14 ", 3); err != nil {
		t.Fatal(err)
	}
	days, _ := svc.Int(context.Background(), model.AdsPendingMaxDays)
	if days != 14 || *repo[model.AdsPendingMaxDays].UpdatedByStaffID != 3 {
		t.Fatalf("days %d, saved %+v", days, repo[model.AdsPendingMaxDays])
	}
}

func TestUpdateRejectsBadValues(t *testing.T) {
	svc := NewSettingService(fakeRepo{})
	for _, value := range []string{"0", "366", "abc", ""} {
		if _, err := svc.Update(context.Background(), model.AdsPendingMaxDays, value, 3); code(t, err) != "SETTING_INVALID_VALUE" {
			t.Errorf("%q: %v", value, err)
		}
	}
	if _, err := svc.Update(context.Background(), "nope", "1", 3); code(t, err) != "SETTING_UNKNOWN_KEY" {
		t.Errorf("unknown key: %v", err)
	}
	if _, err := svc.Update(context.Background(), model.AdsPendingMaxDays, "5", 0); code(t, err) != "SETTING_STAFF_REQUIRED" {
		t.Errorf("no staff: %v", err)
	}
}

func TestCorruptStoredValueFallsBackToDefault(t *testing.T) {
	svc := NewSettingService(fakeRepo{model.AdsPendingMaxDays: {Key: model.AdsPendingMaxDays, Value: "oops"}})
	if days, err := svc.Int(context.Background(), model.AdsPendingMaxDays); err != nil || days != 7 {
		t.Fatalf("got %d, %v", days, err)
	}
}
