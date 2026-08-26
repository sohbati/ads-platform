package model

import "time"

// UserProfile is the 1:1 preference record for a registered user.
type UserProfile struct {
	UserID        int64     `json:"user_id" gorm:"primaryKey"`
	LocationSlugs []string  `json:"location_slugs" gorm:"type:jsonb;serializer:json;not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}
