package repository

import (
	"context"

	"ads-platform/internal/business/setting/model"
)

type SettingRepository interface {
	List(ctx context.Context) ([]model.AppSetting, error)
	Get(ctx context.Context, key string) (*model.AppSetting, error)
	Upsert(ctx context.Context, setting *model.AppSetting) error
}
