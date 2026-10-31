package service

import (
	"context"

	"ads-platform/internal/business/setting/model"
)

type SettingService interface {
	// List returns every known setting, falling back to its default when it was never saved.
	List(ctx context.Context) ([]model.SettingView, error)
	Update(ctx context.Context, key, value string, staffID int64) (*model.SettingView, error)
	// Int reads a known integer setting; background jobs use it.
	Int(ctx context.Context, key string) (int, error)
}
