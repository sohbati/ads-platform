package container

import (
	"ads-bff/internal/business/auth/handler"
	"ads-bff/internal/business/auth/service"
	"ads-bff/internal/core/client/backend"
	"ads-bff/internal/core/client/cache"
	"ads-bff/internal/core/config"
)

type AuthContainer struct {
	AuthHandler *handler.AuthHandler
}

func NewAuthContainer(cfg *config.Config) *AuthContainer {
	backendClient := backend.NewClient(cfg.BackendAPIBaseURL, nil)
	cacheClient := cache.NewClient(cfg.CacheServiceURL, nil)
	authService := service.NewAuthService(cfg, backendClient, cacheClient)

	return &AuthContainer{
		AuthHandler: handler.NewAuthHandler(cfg, authService),
	}
}
