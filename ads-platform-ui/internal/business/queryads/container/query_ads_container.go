package container

import (
	"ads-platform-ui/internal/business/queryads/handler"
	serviceimpl "ads-platform-ui/internal/business/queryads/service/impl"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
)

type QueryAdsContainer struct {
	PageHandler *handler.PageHandler
}

func NewQueryAdsContainer(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog) *QueryAdsContainer {
	svc := serviceimpl.NewQueryAdsService(reg, catalog)
	return &QueryAdsContainer{
		PageHandler: handler.NewPageHandler(svc, cfg),
	}
}
