package container

import (
	"ads-platform-ui/internal/business/newad/handler"
	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
)

type NewAdContainer struct {
	PageHandler *handler.PageHandler
	APIHandler  *handler.APIHandler
}

func NewNewAdContainer(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, cdnClient *cdn.Client, bffClient *bff.Client) *NewAdContainer {
	return &NewAdContainer{
		PageHandler: handler.NewPageHandler(cfg, reg, catalog, cdnClient, bffClient),
		APIHandler:  handler.NewAPIHandler(bffClient),
	}
}
