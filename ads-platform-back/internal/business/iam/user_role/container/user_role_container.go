package container

import (
	rolerepoimpl "ads-platform/internal/business/iam/role/repository/impl"
	userrepoimpl "ads-platform/internal/business/iam/user/repository/impl"
	"ads-platform/internal/business/iam/user_role/handler"
	repoimpl "ads-platform/internal/business/iam/user_role/repository/impl"
	serviceimpl "ads-platform/internal/business/iam/user_role/service/impl"

	"gorm.io/gorm"
)

type UserRoleContainer struct {
	Handler *handler.UserRoleHandler
}

func NewUserRoleContainer(db *gorm.DB) *UserRoleContainer {
	repo := repoimpl.NewUserRoleRepository(db)
	users := userrepoimpl.NewUserRepository(db)
	roles := rolerepoimpl.NewRoleRepository(db)
	svc := serviceimpl.NewUserRoleService(repo, users, roles)
	return &UserRoleContainer{Handler: handler.NewUserRoleHandler(svc)}
}
