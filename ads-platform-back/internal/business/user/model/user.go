package model

import (
	"time"
)

type User struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"type:varchar(100);not null"`
	Mobile     string    `json:"mobile" gorm:"type:varchar(100);not null"`
	NationalId string    `json:"national_id" gorm:"type:varchar(100);not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "user"
}
