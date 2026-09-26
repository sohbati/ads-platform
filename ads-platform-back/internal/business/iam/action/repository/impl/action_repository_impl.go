package impl

import (
	"context"

	"ads-platform/internal/business/iam/action/model"
	"ads-platform/internal/business/iam/action/repository"

	"gorm.io/gorm"
)

type actionRepository struct {
	db *gorm.DB
}

func NewActionRepository(db *gorm.DB) repository.ActionRepository {
	return &actionRepository{db: db}
}

func (r *actionRepository) Create(ctx context.Context, action *model.Action) error {
	return r.db.WithContext(ctx).Create(action).Error
}

func (r *actionRepository) GetByID(ctx context.Context, id int64) (*model.Action, error) {
	var action model.Action
	if err := r.db.WithContext(ctx).First(&action, id).Error; err != nil {
		return nil, err
	}
	return &action, nil
}

func (r *actionRepository) List(ctx context.Context) ([]model.Action, error) {
	var items []model.Action
	err := r.db.WithContext(ctx).Order("code").Find(&items).Error
	if items == nil {
		items = []model.Action{}
	}
	return items, err
}

func (r *actionRepository) Update(ctx context.Context, action *model.Action) error {
	return r.db.WithContext(ctx).Model(action).Select("name").Updates(action).Error
}

func (r *actionRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&model.Action{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
