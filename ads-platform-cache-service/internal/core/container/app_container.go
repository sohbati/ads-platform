package container

import (
	cacheContainer "cache-service/internal/simplecache/otp/container"
)

type AppContainer struct {
	OtpCacheContainer *cacheContainer.OtpCacheContainer
}

func NewAppContainer() *AppContainer {
	return &AppContainer{
		OtpCacheContainer: cacheContainer.NewOtpCacheContainer(),
	}
}
