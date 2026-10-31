package service

import (
	"ads-platform/internal/business/iam/staff/model"
	"context"
)

type StaffInput struct {
	Name       string
	Mobile     string
	NationalId string
	// Password is required on create and ignored on update; use SetPassword to change it.
	Password string
	// IsActive is optional; nil keeps the current value (new staff start active).
	IsActive *bool
}

type StaffPage struct {
	Items    []model.Staff `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"pageSize"`
}

type StaffService interface {
	Create(ctx context.Context, input StaffInput) (*model.Staff, error)
	GetStaffByID(ctx context.Context, id int64) (*model.Staff, error)
	List(ctx context.Context, search string, page, pageSize int) (*StaffPage, error)
	Update(ctx context.Context, id int64, input StaffInput) (*model.Staff, error)
	Delete(ctx context.Context, id int64) error
	SetPassword(ctx context.Context, id int64, password string) error
	// UpsertSuperAdmin creates or promotes the staff member with input.Mobile to an active super admin
	// and sets their password. It is meant for bootstrapping, not for the panel API.
	UpsertSuperAdmin(ctx context.Context, input StaffInput) (*model.Staff, error)
}
