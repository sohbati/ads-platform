package impl

import (
	"context"

	"ads-platform/internal/business/iam/staff_group/model"
	"ads-platform/internal/business/iam/staff_group/repository"

	"gorm.io/gorm"
)

type staffGroupRepository struct {
	db *gorm.DB
}

func NewStaffGroupRepository(db *gorm.DB) repository.StaffGroupRepository {
	return &staffGroupRepository{db: db}
}

func (r *staffGroupRepository) Create(ctx context.Context, group *model.StaffGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *staffGroupRepository) GetByID(ctx context.Context, id int64) (*model.StaffGroup, error) {
	var group model.StaffGroup
	if err := r.db.WithContext(ctx).First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *staffGroupRepository) List(ctx context.Context) ([]model.StaffGroup, error) {
	var items []model.StaffGroup
	err := r.db.WithContext(ctx).Order("code").Find(&items).Error
	if items == nil {
		items = []model.StaffGroup{}
	}
	return items, err
}

func (r *staffGroupRepository) Update(ctx context.Context, group *model.StaffGroup) error {
	return r.db.WithContext(ctx).Model(group).Select("name").Updates(group).Error
}

func (r *staffGroupRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&model.StaffGroup{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *staffGroupRepository) ListStaffIDs(ctx context.Context, groupID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).
		Model(&model.StaffGroupMember{}).
		Where("group_id = ?", groupID).
		Order("staff_id").
		Pluck("staff_id", &ids).Error
	if ids == nil {
		ids = []int64{}
	}
	return ids, err
}

func (r *staffGroupRepository) AddStaff(ctx context.Context, groupID, staffID int64) error {
	return r.db.WithContext(ctx).Create(&model.StaffGroupMember{GroupID: groupID, StaffID: staffID}).Error
}

func (r *staffGroupRepository) RemoveStaff(ctx context.Context, groupID, staffID int64) error {
	res := r.db.WithContext(ctx).
		Where("group_id = ? AND staff_id = ?", groupID, staffID).
		Delete(&model.StaffGroupMember{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *staffGroupRepository) ListByStaff(ctx context.Context, staffID int64) ([]model.StaffGroup, error) {
	var items []model.StaffGroup
	err := r.db.WithContext(ctx).
		Table("staff_group").
		Joins("JOIN staff_group_member ON staff_group_member.group_id = staff_group.id").
		Where("staff_group_member.staff_id = ?", staffID).
		Order("staff_group.code").
		Find(&items).Error
	if items == nil {
		items = []model.StaffGroup{}
	}
	return items, err
}
