package container

import (
	actionrepoimpl "ads-platform/internal/business/iam/action/repository/impl"
	rolerepoimpl "ads-platform/internal/business/iam/role/repository/impl"
	"ads-platform/internal/business/iam/role_action/handler"
	repoimpl "ads-platform/internal/business/iam/role_action/repository/impl"
	serviceimpl "ads-platform/internal/business/iam/role_action/service/impl"

	"gorm.io/gorm"
)

type RoleActionContainer struct {
	Handler *handler.RoleActionHandler
}

func NewRoleActionContainer(db *gorm.DB) *RoleActionContainer {
	repo := repoimpl.NewRoleActionRepository(db)
	roles := rolerepoimpl.NewRoleRepository(db)
	actions := actionrepoimpl.NewActionRepository(db)
	svc := serviceimpl.NewRoleActionService(repo, roles, actions)
	return &RoleActionContainer{Handler: handler.NewRoleActionHandler(svc)}
}
