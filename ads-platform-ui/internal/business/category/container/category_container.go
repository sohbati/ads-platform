package container

import (
	"ads-platform-ui/internal/business/category/handler"
	serviceimpl "ads-platform-ui/internal/business/category/service/impl"
	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
)

type CategoryContainer struct {
	PageHandler *handler.PageHandler
	APIHandler  *handler.APIHandler
}

func NewCategoryContainer(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, cdnClient *cdn.Client) *CategoryContainer {
	svc := serviceimpl.NewCategoryService(cdnClient)
	return &CategoryContainer{
		PageHandler: handler.NewPageHandler(cfg, reg, catalog, svc),
		APIHandler:  handler.NewAPIHandler(svc),
	}
}
