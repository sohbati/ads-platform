package container

import (
	userrepoimpl "ads-platform/internal/business/user/repository/impl"
	"ads-platform/internal/business/userprofile/handler"
	repoimpl "ads-platform/internal/business/userprofile/repository/impl"
	serviceimpl "ads-platform/internal/business/userprofile/service/impl"

	"gorm.io/gorm"
)

type UserProfileContainer struct {
	Handler *handler.UserProfileHandler
}

func NewUserProfileContainer(db *gorm.DB) *UserProfileContainer {
	profiles := repoimpl.NewUserProfileRepository(db)
	users := userrepoimpl.NewUserRepository(db)
	svc := serviceimpl.NewUserProfileService(profiles, users)
	return &UserProfileContainer{
		Handler: handler.NewUserProfileHandler(svc),
	}
}
