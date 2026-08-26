package container

import (
	"ads-bff/internal/business/ads/handler"
	"ads-bff/internal/business/auth/service"
	backendclient "ads-bff/internal/core/client/backend"
	"ads-bff/internal/core/config"
)

type AdsContainer struct {
	AdHandler *handler.AdHandler
}

func NewAdsContainer(cfg *config.Config, auth service.AuthService) *AdsContainer {
	backend := backendclient.NewAdsClient(cfg.BackendAPIBaseURL)
	return &AdsContainer{
		AdHandler: handler.NewAdHandler(cfg, auth, backend),
	}
}
