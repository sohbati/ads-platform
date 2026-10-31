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

func (r *staffRepository) ListStaff(ctx context.Context, search []string, limit, offset int) ([]model.Staff, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Staff{})
	if len(search) > 0 {
		cond := r.db
		for _, term := range search {
			like := "%" + term + "%"
			cond = cond.Or("name ILIKE ? OR mobile LIKE ? OR national_id LIKE ?", like, like, like)
		}
		query = query.Where(cond)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := []model.Staff{}
	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&items).Error
	return items, total, err
}

func (r *staffRepository) UpdateStaff(ctx context.Context, staff *model.Staff) error {
	return r.db.WithContext(ctx).Model(staff).Select("name", "mobile", "national_id", "password_hash", "is_super_admin", "is_active").Updates(staff).Error
}

func (r *staffRepository) DeleteStaff(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&model.Staff{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *staffRepository) GetStaffByMobile(ctx context.Context, mobile string) (*model.Staff, error) {
	var staff model.Staff
	err := r.db.WithContext(ctx).Where("mobile = ?", mobile).First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}
