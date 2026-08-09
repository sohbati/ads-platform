package container

import (
	"ads-platform-ui/internal/business/myinfo/handler"
	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
)

type MyInfoContainer struct {
	PageHandler *handler.PageHandler
}

func NewMyInfoContainer(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, bffClient *bff.Client) *MyInfoContainer {
	return &MyInfoContainer{
		PageHandler: handler.NewPageHandler(cfg, reg, catalog, bffClient),
	}
}
