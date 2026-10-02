package container

import (
	staffrepoimpl "ads-platform/internal/business/iam/staff/repository/impl"
	"ads-platform/internal/business/iam/staff_group/handler"
	repoimpl "ads-platform/internal/business/iam/staff_group/repository/impl"
	serviceimpl "ads-platform/internal/business/iam/staff_group/service/impl"

	"gorm.io/gorm"
)

type StaffGroupContainer struct {
	Handler *handler.StaffGroupHandler
}

func NewStaffGroupContainer(db *gorm.DB) *StaffGroupContainer {
	repo := repoimpl.NewStaffGroupRepository(db)
	staff := staffrepoimpl.NewStaffRepository(db)
	svc := serviceimpl.NewStaffGroupService(repo, staff)
	return &StaffGroupContainer{Handler: handler.NewStaffGroupHandler(svc)}
}
