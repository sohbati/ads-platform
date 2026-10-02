package container

import (
	"ads-platform/internal/business/iam/staff/handler"
	repoimpl "ads-platform/internal/business/iam/staff/repository/impl"
	serviceimpl "ads-platform/internal/business/iam/staff/service/impl"

	"gorm.io/gorm"
)

type StaffContainer struct {
	Handler *handler.StaffHandler
}

func NewStaffContainer(db *gorm.DB) *StaffContainer {
	staffRepo := repoimpl.NewStaffRepository(db)
	staffService := serviceimpl.NewStaffService(staffRepo)
	return &StaffContainer{
		Handler: handler.NewStaffHandler(staffService),
	}
}
