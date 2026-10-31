package impl

import (
	"context"

	"ads-platform/internal/business/iam/role/model"
	"ads-platform/internal/business/iam/role/repository"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id int64) (*model.Role, error) {
	var role model.Role
	if err := r.db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) List(ctx context.Context) ([]model.Role, error) {
	var items []model.Role
	err := r.db.WithContext(ctx).Order("code").Find(&items).Error
	if items == nil {
		items = []model.Role{}
	}
	return items, err
}

func (r *roleRepository) Update(ctx context.Context, role *model.Role) error {
	return r.db.WithContext(ctx).Model(role).Select("code", "name").Updates(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&model.Role{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
