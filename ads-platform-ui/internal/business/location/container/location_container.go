package container

import (
	"ads-platform-ui/internal/business/location/handler"
	serviceimpl "ads-platform-ui/internal/business/location/service/impl"
	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
)

type LocationContainer struct {
	PageHandler *handler.PageHandler
	APIHandler  *handler.APIHandler
}

func NewLocationContainer(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, cdnClient *cdn.Client) *LocationContainer {
	svc := serviceimpl.NewLocationService(cdnClient)
	return &LocationContainer{
		PageHandler: handler.NewPageHandler(cfg, reg, catalog),
		APIHandler:  handler.NewAPIHandler(svc),
	}
}
