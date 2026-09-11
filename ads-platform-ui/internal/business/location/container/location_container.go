package container

import (
	"ads-platform-ui/internal/business/location/handler"
	serviceimpl "ads-platform-ui/internal/business/location/service/impl"
	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/config"
)

type LocationContainer struct {
	APIHandler *handler.APIHandler
	GeoHandler *handler.GeoHandler
}

func NewLocationContainer(cdnClient *cdn.Client, cfg *config.Config) *LocationContainer {
	svc := serviceimpl.NewLocationService(cdnClient)
	base := ""
	if cfg != nil {
		base = cfg.NominatimBaseURL
	}
	return &LocationContainer{
		APIHandler: handler.NewAPIHandler(svc),
		GeoHandler: handler.NewGeoHandler(serviceimpl.NewNominatimClient(base)),
	}
}
