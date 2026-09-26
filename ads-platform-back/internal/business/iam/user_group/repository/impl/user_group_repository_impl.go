package impl

import (
	"context"

	"ads-platform/internal/business/iam/user_group/model"
	"ads-platform/internal/business/iam/user_group/repository"

	"gorm.io/gorm"
)

type userGroupRepository struct {
	db *gorm.DB
}

func NewUserGroupRepository(db *gorm.DB) repository.UserGroupRepository {
	return &userGroupRepository{db: db}
}

func (r *userGroupRepository) Create(ctx context.Context, group *model.UserGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *userGroupRepository) GetByID(ctx context.Context, id int64) (*model.UserGroup, error) {
	var group model.UserGroup
	if err := r.db.WithContext(ctx).First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *userGroupRepository) List(ctx context.Context) ([]model.UserGroup, error) {
	var items []model.UserGroup
	err := r.db.WithContext(ctx).Order("code").Find(&items).Error
	if items == nil {
		items = []model.UserGroup{}
	}
	return items, err
}

func (r *userGroupRepository) Update(ctx context.Context, group *model.UserGroup) error {
	return r.db.WithContext(ctx).Model(group).Select("name").Updates(group).Error
}

func (r *userGroupRepository) Delete(ctx context.Context, id int64) error {
	res := r.db.WithContext(ctx).Delete(&model.UserGroup{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userGroupRepository) ListUserIDs(ctx context.Context, groupID int64) ([]int64, error) {
	var ids []int64
	err := r.db.WithContext(ctx).
		Model(&model.UserGroupMember{}).
		Where("group_id = ?", groupID).
		Order("user_id").
		Pluck("user_id", &ids).Error
	if ids == nil {
		ids = []int64{}
	}
	return ids, err
}

func (r *userGroupRepository) AddUser(ctx context.Context, groupID, userID int64) error {
	return r.db.WithContext(ctx).Create(&model.UserGroupMember{GroupID: groupID, UserID: userID}).Error
}

func (r *userGroupRepository) RemoveUser(ctx context.Context, groupID, userID int64) error {
	res := r.db.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&model.UserGroupMember{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userGroupRepository) ListByUser(ctx context.Context, userID int64) ([]model.UserGroup, error) {
	var items []model.UserGroup
	err := r.db.WithContext(ctx).
		Table("user_group").
		Joins("JOIN user_group_member ON user_group_member.group_id = user_group.id").
		Where("user_group_member.user_id = ?", userID).
		Order("user_group.code").
		Find(&items).Error
	if items == nil {
		items = []model.UserGroup{}
	}
	return items, err
}
