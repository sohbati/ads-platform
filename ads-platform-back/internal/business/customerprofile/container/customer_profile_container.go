package container

import (
	customerrepoimpl "ads-platform/internal/business/customer/repository/impl"
	"ads-platform/internal/business/customerprofile/handler"
	repoimpl "ads-platform/internal/business/customerprofile/repository/impl"
	serviceimpl "ads-platform/internal/business/customerprofile/service/impl"

	"gorm.io/gorm"
)

type CustomerProfileContainer struct {
	Handler *handler.CustomerProfileHandler
}

func NewCustomerProfileContainer(db *gorm.DB) *CustomerProfileContainer {
	profiles := repoimpl.NewCustomerProfileRepository(db)
	customers := customerrepoimpl.NewCustomerRepository(db)
	svc := serviceimpl.NewCustomerProfileService(profiles, customers)
	return &CustomerProfileContainer{
		Handler: handler.NewCustomerProfileHandler(svc),
	}
}
