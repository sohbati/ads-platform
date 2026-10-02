package service

import (
	"ads-platform/internal/business/customer/model"
	"context"
)

type CustomerService interface {
	GetCustomerByID(ctx context.Context, id int64) (*model.Customer, error)
	GetCustomerByMobile(ctx context.Context, mobile string) (*model.Customer, error)
	RegisterByMobile(ctx context.Context, mobile string) (*model.Customer, error)
}
