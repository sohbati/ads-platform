package repository

import (
	"ads-platform/internal/business/customer/model"
	"context"
)

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *model.Customer) error
	GetCustomerByID(ctx context.Context, id int64) (*model.Customer, error)
	GetCustomerByMobile(ctx context.Context, mobile string) (*model.Customer, error)
}
