package container

import (
	staffrepoimpl "ads-platform/internal/business/iam/staff/repository/impl"
	"ads-platform/internal/business/iam/staff_auth/handler"
	serviceimpl "ads-platform/internal/business/iam/staff_auth/service/impl"
	staffrolerepoimpl "ads-platform/internal/business/iam/staff_role/repository/impl"
	"ads-platform/internal/core/mobile"

	"gorm.io/gorm"
)

type StaffAuthContainer struct {
	Handler *handler.StaffAuthHandler
}

func NewStaffAuthContainer(db *gorm.DB, defaultCountryCode string) *StaffAuthContainer {
	svc := serviceimpl.NewStaffAuthService(
		staffrepoimpl.NewStaffRepository(db),
		staffrolerepoimpl.NewStaffRoleRepository(db),
		mobile.NewNormalizer(defaultCountryCode),
	)
	return &StaffAuthContainer{Handler: handler.NewStaffAuthHandler(svc)}
}
