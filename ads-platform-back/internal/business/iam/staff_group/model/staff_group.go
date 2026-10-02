package model

import "time"

type StaffGroup struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"type:varchar(80);uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"type:varchar(120);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (StaffGroup) TableName() string { return "staff_group" }

type StaffGroupMember struct {
	GroupID int64 `json:"group_id" gorm:"primaryKey"`
	StaffID  int64 `json:"staff_id" gorm:"primaryKey"`
}

func (StaffGroupMember) TableName() string { return "staff_group_member" }
