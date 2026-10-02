package model

type StaffRole struct {
	StaffID int64 `json:"staff_id" gorm:"primaryKey"`
	RoleID int64 `json:"role_id" gorm:"primaryKey"`
}

func (StaffRole) TableName() string { return "staff_role" }
