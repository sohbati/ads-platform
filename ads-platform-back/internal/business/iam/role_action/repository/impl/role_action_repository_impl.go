package impl

import (
	"context"

	actionmodel "ads-platform/internal/business/iam/action/model"
	"ads-platform/internal/business/iam/role_action/model"
	"ads-platform/internal/business/iam/role_action/repository"

	"gorm.io/gorm"
)

type roleActionRepository struct {
	db *gorm.DB
}

func NewRoleActionRepository(db *gorm.DB) repository.RoleActionRepository {
	return &roleActionRepository{db: db}
}

func (r *roleActionRepository) ListActions(ctx context.Context, roleID int64) ([]actionmodel.Action, error) {
	var items []actionmodel.Action
	err := r.db.WithContext(ctx).
		Table("action").
		Joins("JOIN role_action ON role_action.action_id = action.id").
		Where("role_action.role_id = ?", roleID).
		Order("action.code").
		Find(&items).Error
	if items == nil {
		items = []actionmodel.Action{}
	}
	return items, err
}

func (r *roleActionRepository) Replace(ctx context.Context, roleID int64, actionIDs []int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleAction{}).Error; err != nil {
			return err
		}
		for _, actionID := range actionIDs {
			if err := tx.Create(&model.RoleAction{RoleID: roleID, ActionID: actionID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleActionRepository) Add(ctx context.Context, roleID, actionID int64) error {
	return r.db.WithContext(ctx).Create(&model.RoleAction{RoleID: roleID, ActionID: actionID}).Error
}

func (r *roleActionRepository) Remove(ctx context.Context, roleID, actionID int64) error {
	res := r.db.WithContext(ctx).
		Where("role_id = ? AND action_id = ?", roleID, actionID).
		Delete(&model.RoleAction{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
