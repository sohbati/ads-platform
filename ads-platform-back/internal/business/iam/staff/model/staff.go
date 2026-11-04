package model

import (
	"time"
)

type Staff struct {
	ID           int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string `json:"name" gorm:"type:varchar(100);not null"`
	Mobile       string `json:"mobile" gorm:"type:varchar(100);not null"`
	NationalId   string `json:"national_id" gorm:"type:varchar(100);not null"`
	PasswordHash string `json:"-" gorm:"type:varchar(255);not null"`
	IsSuperAdmin bool   `json:"is_super_admin" gorm:"not null"`
	IsActive     bool   `json:"is_active" gorm:"not null"`
	// IsReviewer: one of the staff member's roles grants advertisement.review. Filled by ListStaff only.
	IsReviewer bool      `json:"is_reviewer" gorm:"->;-:migration"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (Staff) TableName() string {
	return "staff"
}
