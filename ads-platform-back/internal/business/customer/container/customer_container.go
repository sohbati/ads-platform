package container

import (
	"ads-platform/internal/business/customer/handler"
	repoimpl "ads-platform/internal/business/customer/repository/impl"
	serviceimpl "ads-platform/internal/business/customer/service/impl"
	"ads-platform/internal/core/mobile"

	"gorm.io/gorm"
)

type CustomerContainer struct {
	Handler *handler.CustomerHandler
}

func NewCustomerContainer(db *gorm.DB, defaultCountryCode string) *CustomerContainer {
	customers := repoimpl.NewCustomerRepository(db)
	svc := serviceimpl.NewCustomerService(customers)
	mobileNorm := mobile.NewNormalizer(defaultCountryCode)
	return &CustomerContainer{
		Handler: handler.NewCustomerHandler(svc, mobileNorm),
	}
}
