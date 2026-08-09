package container

import (
	cacheContainer "cache-service/internal/simplecache/otp/container"
	sessionContainer "cache-service/internal/simplecache/session/container"
)

type AppContainer struct {
	OtpCacheContainer     *cacheContainer.OtpCacheContainer
	SessionCacheContainer *sessionContainer.SessionCacheContainer
}

func NewAppContainer() *AppContainer {
	return &AppContainer{
		OtpCacheContainer:     cacheContainer.NewOtpCacheContainer(),
		SessionCacheContainer: sessionContainer.NewSessionCacheContainer(),
	}
}
