package impl

import (
	"context"

	actionmodel "ads-platform/internal/business/iam/action/model"
	rolemodel "ads-platform/internal/business/iam/role/model"
	"ads-platform/internal/business/iam/user_role/model"
	"ads-platform/internal/business/iam/user_role/repository"

	"gorm.io/gorm"
)

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) repository.UserRoleRepository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) ListRoles(ctx context.Context, userID int64) ([]rolemodel.Role, error) {
	var items []rolemodel.Role
	err := r.db.WithContext(ctx).
		Table(`"role"`).
		Joins("JOIN user_role ON user_role.role_id = \"role\".id").
		Where("user_role.user_id = ?", userID).
		Order(`"role".code`).
		Find(&items).Error
	if items == nil {
		items = []rolemodel.Role{}
	}
	return items, err
}

func (r *userRoleRepository) Add(ctx context.Context, userID, roleID int64) error {
	return r.db.WithContext(ctx).Create(&model.UserRole{UserID: userID, RoleID: roleID}).Error
}

func (r *userRoleRepository) Remove(ctx context.Context, userID, roleID int64) error {
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&model.UserRole{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRoleRepository) ListActions(ctx context.Context, userID int64) ([]actionmodel.Action, error) {
	var items []actionmodel.Action
	err := r.db.WithContext(ctx).
		Table("action").
		Distinct("action.id", "action.code", "action.name", "action.created_at", "action.updated_at").
		Joins("JOIN role_action ON role_action.action_id = action.id").
		Joins("JOIN user_role ON user_role.role_id = role_action.role_id").
		Where("user_role.user_id = ?", userID).
		Order("action.code").
		Find(&items).Error
	if items == nil {
		items = []actionmodel.Action{}
	}
	return items, err
}

func (r *userRoleRepository) HasAction(ctx context.Context, userID int64, code string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).
		Table("action").
		Joins("JOIN role_action ON role_action.action_id = action.id").
		Joins("JOIN user_role ON user_role.role_id = role_action.role_id").
		Where("user_role.user_id = ? AND action.code = ?", userID, code).
		Count(&n).Error
	return n > 0, err
}
