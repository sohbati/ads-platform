package model

import "time"

// CustomerProfile is the 1:1 preference record for a registered customer.
type CustomerProfile struct {
	CustomerID        int64     `json:"customer_id" gorm:"primaryKey"`
	LocationSlugs []string  `json:"location_slugs" gorm:"type:jsonb;serializer:json;not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (CustomerProfile) TableName() string {
	return "customer_profile"
}
