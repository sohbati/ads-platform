package model

import "time"

type UserGroup struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"type:varchar(80);uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"type:varchar(120);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserGroup) TableName() string { return "user_group" }

type UserGroupMember struct {
	GroupID int64 `json:"group_id" gorm:"primaryKey"`
	UserID  int64 `json:"user_id" gorm:"primaryKey"`
}

func (UserGroupMember) TableName() string { return "user_group_member" }
