package container

// Container acts as the application's dependency container (composition root).
// Its responsibility is to initialize and wire together all core components
// such as repositories, services, and handlers.
//
// The container centralizes dependency creation so that other parts of the
// application (like the router) do not need to know how objects are constructed.
// Instead, they simply receive ready-to-use instances.
//
// This approach provides several benefits:
// - Keeps object creation in one place
// - Promotes clean separation of concerns
// - Makes testing easier by allowing dependencies to be replaced with mocks
// - Prevents tight coupling between layers
//
// In general, the dependency flow in the application is:
//
//   Handler → Service → Repository
//
// The container is responsible for constructing these layers and injecting
// them into each other during application startup.

import (
	"log"
	"time"

	cachestore "cache-service/internal/cachestore"
	"cache-service/internal/simplecache/otp/handler"
	serviceimpl "cache-service/internal/simplecache/otp/service/impl"
)

type OtpCacheContainer struct {
	OtpCacheHandler *handler.OtpCacheHandler
	OtpCacheStore   *cachestore.Cache[string, string]
}

func NewOtpCacheContainer() *OtpCacheContainer {

	// Create the store with full config
	store := cachestore.New[string, string](
		cachestore.WithTTL(500*time.Second),
		cachestore.WithCleanupInterval(30*time.Second),
		cachestore.WithOnEvict(func(key string) {
			log.Println("  [evicted]:", key)
		}),
	)
	// services
	cacheService := serviceimpl.NewOtpCacheService(store)

	// handlers
	cacheHandler := handler.NewOtpCacheHandler(cacheService)
	log.Println("Otp cache configured")
	return &OtpCacheContainer{
		OtpCacheHandler: cacheHandler,
		OtpCacheStore:   store,
	}
}
