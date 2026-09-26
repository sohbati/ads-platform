package container

import (
	"ads-platform/internal/business/iam/role/handler"
	repoimpl "ads-platform/internal/business/iam/role/repository/impl"
	serviceimpl "ads-platform/internal/business/iam/role/service/impl"

	"gorm.io/gorm"
)

type RoleContainer struct {
	Handler *handler.RoleHandler
}

func NewRoleContainer(db *gorm.DB) *RoleContainer {
	repo := repoimpl.NewRoleRepository(db)
	svc := serviceimpl.NewRoleService(repo)
	return &RoleContainer{Handler: handler.NewRoleHandler(svc)}
}
