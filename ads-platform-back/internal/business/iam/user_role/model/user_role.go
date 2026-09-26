package model

type UserRole struct {
	UserID int64 `json:"user_id" gorm:"primaryKey"`
	RoleID int64 `json:"role_id" gorm:"primaryKey"`
}

func (UserRole) TableName() string { return "user_role" }
