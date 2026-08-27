package container

import (
	"ads-platform-ui/internal/business/location/handler"
	serviceimpl "ads-platform-ui/internal/business/location/service/impl"
	"ads-platform-ui/internal/core/cdn"
)

type LocationContainer struct {
	APIHandler *handler.APIHandler
}

func NewLocationContainer(cdnClient *cdn.Client) *LocationContainer {
	svc := serviceimpl.NewLocationService(cdnClient)
	return &LocationContainer{
		APIHandler: handler.NewAPIHandler(svc),
	}
}
