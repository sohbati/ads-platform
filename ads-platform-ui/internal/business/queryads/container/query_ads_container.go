package container

import (
	searchclient "ads-platform-ui/internal/business/queryads/client"
	"ads-platform-ui/internal/business/queryads/handler"
	serviceimpl "ads-platform-ui/internal/business/queryads/service/impl"
	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
)

type QueryAdsContainer struct {
	PageHandler *handler.PageHandler
}

func NewQueryAdsContainer(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, bffClient *bff.Client) *QueryAdsContainer {
	svc := serviceimpl.NewQueryAdsService(reg, catalog, searchclient.NewSearchClient(bffClient))
	return &QueryAdsContainer{
		PageHandler: handler.NewPageHandler(svc, cfg),
	}
}
