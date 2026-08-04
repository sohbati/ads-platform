package container

import (
	"gorm.io/gorm"
	userContainer "ads-platform/internal/business/user/container"

)

type AppContainer struct {
	User *userContainer.UserContainer
}

func NewAppContainer(db *gorm.DB) *AppContainer {
	return &AppContainer{
		User: userContainer.NewUserContainer(db),
	}
}
