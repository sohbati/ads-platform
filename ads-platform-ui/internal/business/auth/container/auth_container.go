package container

import (
	"ads-platform-ui/internal/business/auth/handler"
	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
)

type AuthContainer struct {
	APIHandler  *handler.APIHandler
	PageHandler *handler.PageHandler
}

func NewAuthContainer(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, client *bff.Client) *AuthContainer {
	return &AuthContainer{
		APIHandler:  handler.NewAPIHandler(client),
		PageHandler: handler.NewPageHandler(cfg, reg, catalog),
	}
}
