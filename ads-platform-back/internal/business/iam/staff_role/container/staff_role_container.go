package container

import (
	rolerepoimpl "ads-platform/internal/business/iam/role/repository/impl"
	staffrepoimpl "ads-platform/internal/business/iam/staff/repository/impl"
	"ads-platform/internal/business/iam/staff_role/handler"
	repoimpl "ads-platform/internal/business/iam/staff_role/repository/impl"
	serviceimpl "ads-platform/internal/business/iam/staff_role/service/impl"

	"gorm.io/gorm"
)

type StaffRoleContainer struct {
	Handler *handler.StaffRoleHandler
}

func NewStaffRoleContainer(db *gorm.DB) *StaffRoleContainer {
	repo := repoimpl.NewStaffRoleRepository(db)
	staff := staffrepoimpl.NewStaffRepository(db)
	roles := rolerepoimpl.NewRoleRepository(db)
	svc := serviceimpl.NewStaffRoleService(repo, staff, roles)
	return &StaffRoleContainer{Handler: handler.NewStaffRoleHandler(svc)}
}
