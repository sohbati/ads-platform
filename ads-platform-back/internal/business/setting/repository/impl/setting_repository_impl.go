package impl

import (
	"context"
	"errors"

	"ads-platform/internal/business/setting/model"
	"ads-platform/internal/business/setting/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) repository.SettingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) List(ctx context.Context) ([]model.AppSetting, error) {
	var items []model.AppSetting
	err := r.db.WithContext(ctx).Order("key").Find(&items).Error
	return items, err
}

func (r *settingRepository) Get(ctx context.Context, key string) (*model.AppSetting, error) {
	var setting model.AppSetting
	err := r.db.WithContext(ctx).Where("key = ?", key).Take(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *settingRepository) Upsert(ctx context.Context, setting *model.AppSetting) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_by_staff_id", "updated_at"}),
	}).Create(setting).Error
}
