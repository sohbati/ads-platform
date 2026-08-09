package container

import (
	"log"
	"time"

	cachestore "cache-service/internal/cachestore"
	"cache-service/internal/simplecache/session/handler"
	serviceimpl "cache-service/internal/simplecache/session/service/impl"
)

type SessionCacheContainer struct {
	SessionCacheHandler *handler.SessionCacheHandler
	SessionCacheStore   *cachestore.Cache[string, string]
}

func NewSessionCacheContainer() *SessionCacheContainer {
	store := cachestore.New[string, string](
		cachestore.WithTTL(24*time.Hour),
		cachestore.WithCleanupInterval(30*time.Second),
		cachestore.WithOnEvict(func(key string) {
			log.Println("  [session evicted]:", key)
		}),
	)

	cacheService := serviceimpl.NewSessionCacheService(store)
	cacheHandler := handler.NewSessionCacheHandler(cacheService)
	log.Println("Session cache configured")

	return &SessionCacheContainer{
		SessionCacheHandler: cacheHandler,
		SessionCacheStore:   store,
	}
}
