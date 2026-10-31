package service

import (
	"context"

	staffmodel "ads-platform/internal/business/iam/staff/model"
)

// SuperAdminPermission grants every panel action.
const SuperAdminPermission = "*"

type StaffSession struct {
	Staff       staffmodel.Staff `json:"staff"`
	Permissions []string         `json:"permissions"`
}

type StaffAuthService interface {
	Login(ctx context.Context, mobile, password string) (*StaffSession, error)
	// Session reloads a signed-in staff member, so deactivation and role changes apply immediately.
	Session(ctx context.Context, staffID int64) (*StaffSession, error)
}
