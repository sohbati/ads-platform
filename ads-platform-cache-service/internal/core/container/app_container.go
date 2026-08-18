package container

import (
	"cache-service/internal/core/config"
	cityContainer "cache-service/internal/simplecache/city/container"
	cacheContainer "cache-service/internal/simplecache/otp/container"
	sessionContainer "cache-service/internal/simplecache/session/container"
)

type AppContainer struct {
	OtpCacheContainer     *cacheContainer.OtpCacheContainer
	SessionCacheContainer *sessionContainer.SessionCacheContainer
	CityCacheContainer    *cityContainer.CityCacheContainer
}

func NewAppContainer(cfg *config.Config) *AppContainer {
	return &AppContainer{
		OtpCacheContainer:     cacheContainer.NewOtpCacheContainer(),
		SessionCacheContainer: sessionContainer.NewSessionCacheContainer(),
		CityCacheContainer:    cityContainer.NewCityCacheContainer(cfg.CDNBaseURL),
	}
}
