package container

import (
	userrepoimpl "ads-platform/internal/business/iam/user/repository/impl"
	"ads-platform/internal/business/iam/user_group/handler"
	repoimpl "ads-platform/internal/business/iam/user_group/repository/impl"
	serviceimpl "ads-platform/internal/business/iam/user_group/service/impl"

	"gorm.io/gorm"
)

type UserGroupContainer struct {
	Handler *handler.UserGroupHandler
}

func NewUserGroupContainer(db *gorm.DB) *UserGroupContainer {
	repo := repoimpl.NewUserGroupRepository(db)
	users := userrepoimpl.NewUserRepository(db)
	svc := serviceimpl.NewUserGroupService(repo, users)
	return &UserGroupContainer{Handler: handler.NewUserGroupHandler(svc)}
}
