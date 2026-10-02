package impl

import (
	"ads-platform/internal/business/customer/model"
	"ads-platform/internal/business/customer/repository"
	"context"

	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) repository.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepository) GetCustomerByID(ctx context.Context, id int64) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.WithContext(ctx).First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetCustomerByMobile(ctx context.Context, mobile string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.WithContext(ctx).Where("mobile = ?", mobile).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}
