package impl

import (
	"context"

	actionmodel "ads-platform/internal/business/iam/action/model"
	rolemodel "ads-platform/internal/business/iam/role/model"
	"ads-platform/internal/business/iam/staff_role/model"
	"ads-platform/internal/business/iam/staff_role/repository"

	"gorm.io/gorm"
)

type staffRoleRepository struct {
	db *gorm.DB
}

func NewStaffRoleRepository(db *gorm.DB) repository.StaffRoleRepository {
	return &staffRoleRepository{db: db}
}

func (r *staffRoleRepository) ListRoles(ctx context.Context, staffID int64) ([]rolemodel.Role, error) {
	var items []rolemodel.Role
	err := r.db.WithContext(ctx).
		Table(`"role"`).
		Joins("JOIN staff_role ON staff_role.role_id = \"role\".id").
		Where("staff_role.staff_id = ?", staffID).
		Order(`"role".code`).
		Find(&items).Error
	if items == nil {
		items = []rolemodel.Role{}
	}
	return items, err
}

func (r *staffRoleRepository) Add(ctx context.Context, staffID, roleID int64) error {
	return r.db.WithContext(ctx).Create(&model.StaffRole{StaffID: staffID, RoleID: roleID}).Error
}

func (r *staffRoleRepository) Remove(ctx context.Context, staffID, roleID int64) error {
	res := r.db.WithContext(ctx).
		Where("staff_id = ? AND role_id = ?", staffID, roleID).
		Delete(&model.StaffRole{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *staffRoleRepository) ListActions(ctx context.Context, staffID int64) ([]actionmodel.Action, error) {
	var items []actionmodel.Action
	err := r.db.WithContext(ctx).
		Table("action").
		Distinct("action.id", "action.code", "action.name", "action.created_at", "action.updated_at").
		Joins("JOIN role_action ON role_action.action_id = action.id").
		Joins("JOIN staff_role ON staff_role.role_id = role_action.role_id").
		Where("staff_role.staff_id = ?", staffID).
		Order("action.code").
		Find(&items).Error
	if items == nil {
		items = []actionmodel.Action{}
	}
	return items, err
}

func (r *staffRoleRepository) HasAction(ctx context.Context, staffID int64, code string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).
		Table("action").
		Joins("JOIN role_action ON role_action.action_id = action.id").
		Joins("JOIN staff_role ON staff_role.role_id = role_action.role_id").
		Where("staff_role.staff_id = ? AND action.code = ?", staffID, code).
		Count(&n).Error
	return n > 0, err
}
