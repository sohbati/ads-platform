package model

import "time"

type Action struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"type:varchar(80);uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"type:varchar(120);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Action) TableName() string { return "action" }
