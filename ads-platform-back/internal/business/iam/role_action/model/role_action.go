package model

type RoleAction struct {
	RoleID   int64 `json:"role_id" gorm:"primaryKey"`
	ActionID int64 `json:"action_id" gorm:"primaryKey"`
}

func (RoleAction) TableName() string { return "role_action" }
