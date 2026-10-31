package model

import "time"

// AppSetting maps to ads_platform_schema.app_setting.
type AppSetting struct {
	Key              string    `json:"key" gorm:"primaryKey;type:varchar(100)"`
	Value            string    `json:"value" gorm:"type:text;not null"`
	UpdatedByStaffID *int64    `json:"updated_by_staff_id"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (AppSetting) TableName() string { return "app_setting" }

// Known setting keys. Add new keys to the registry in the service as well.
const (
	// AdsPendingMaxDays is how many days an ad may stay pending before the timeout job acts on it.
	AdsPendingMaxDays = "ads.pending_max_days"
)

// SettingView is one known setting with its current value (the default when never saved).
type SettingView struct {
	Key              string     `json:"key"`
	Value            string     `json:"value"`
	Type             string     `json:"type"`
	Min              *int       `json:"min,omitempty"`
	Max              *int       `json:"max,omitempty"`
	Default          string     `json:"default"`
	UpdatedByStaffID *int64     `json:"updated_by_staff_id"`
	UpdatedAt        *time.Time `json:"updated_at"`
}
