package impl

import (
	"ads-platform/internal/business/customer/errorcode"
	"ads-platform/internal/business/customer/model"
	"ads-platform/internal/business/customer/repository"
	"ads-platform/internal/business/customer/service"
	"ads-platform/internal/core/exception"
	"context"
	"errors"

	"gorm.io/gorm"
)

type customerService struct {
	customers repository.CustomerRepository
}

func NewCustomerService(customers repository.CustomerRepository) service.CustomerService {
	return &customerService{customers: customers}
}

func (s *customerService) GetCustomerByID(ctx context.Context, id int64) (*model.Customer, error) {
	if id <= 0 {
		return nil, exception.NewAppError(
			errorcode.ErrCustomerNotFound.Code, errorcode.ErrCustomerNotFound.HttpStatus)
	}
	customer, err := s.customers.GetCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(
				errorcode.ErrCustomerNotFound.Code, errorcode.ErrCustomerNotFound.HttpStatus).WithCause(err)
		}
		return nil, err
	}
	return customer, nil
}

func (s *customerService) GetCustomerByMobile(ctx context.Context, mobile string) (*model.Customer, error) {
	if mobile == "" {
		return nil, exception.NewAppError(
			errorcode.ErrMobileEmpty.Code, errorcode.ErrMobileEmpty.HttpStatus, mobile)
	}
	customer, err := s.customers.GetCustomerByMobile(ctx, mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(
				errorcode.ErrCustomerNotFound.Code, errorcode.ErrCustomerNotFound.HttpStatus, mobile).WithCause(err)
		}
		return nil, err
	}
	return customer, nil
}

func (s *customerService) RegisterByMobile(ctx context.Context, mobile string) (*model.Customer, error) {
	if mobile == "" {
		return nil, exception.NewAppError(
			errorcode.ErrMobileEmpty.Code, errorcode.ErrMobileEmpty.HttpStatus, mobile)
	}

	existing, err := s.customers.GetCustomerByMobile(ctx, mobile)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	customer := &model.Customer{
		Mobile:     mobile,
		Name:       mobile,
		NationalId: "",
	}
	if err := s.customers.CreateCustomer(ctx, customer); err != nil {
		existing, getErr := s.customers.GetCustomerByMobile(ctx, mobile)
		if getErr == nil {
			return existing, nil
		}
		return nil, err
	}
	return customer, nil
}
