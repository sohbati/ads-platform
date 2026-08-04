package container

import (
	"ads-platform-ui/internal/business/newad/handler"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
)

type NewAdContainer struct {
	PageHandler *handler.PageHandler
}

func NewNewAdContainer(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog) *NewAdContainer {
	return &NewAdContainer{
		PageHandler: handler.NewPageHandler(cfg, reg, catalog),
	}
}
