package impl

import (
	"ads-platform/internal/business/iam/staff/model"
	"ads-platform/internal/business/iam/staff/repository"
	"context"

	"gorm.io/gorm"
)

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) repository.StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) CreateStaff(ctx context.Context, staff *model.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

func (r *staffRepository) GetStaffByID(ctx context.Context, id int64) (*model.Staff, error) {
	var staff model.Staff
	err := r.db.WithContext(ctx).First(&staff, id).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) GetStaffList(ctx context.Context, limit, offset int) ([]*model.Staff, error) {
	var items []*model.Staff
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&items).Error
	return items, err
}

func (r *staffRepository) UpdateStaff(ctx context.Context, staff *model.Staff) error {
	return r.db.WithContext(ctx).Save(staff).Error
}

func (r *staffRepository) DeleteStaff(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Staff{}, id).Error
}

func (r *staffRepository) GetStaffByMobile(ctx context.Context, mobile string) (*model.Staff, error) {
	var staff model.Staff
	err := r.db.WithContext(ctx).Where("mobile = ?", mobile).First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}
