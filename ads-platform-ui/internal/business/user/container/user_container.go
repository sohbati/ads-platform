package container

import (
	"ads-platform-ui/internal/business/user/handler"
	"ads-platform-ui/internal/core/bff"
)

type UserContainer struct {
	APIHandler *handler.APIHandler
}

func NewUserContainer(client *bff.Client) *UserContainer {
	return &UserContainer{
		APIHandler: handler.NewAPIHandler(client),
	}
}
