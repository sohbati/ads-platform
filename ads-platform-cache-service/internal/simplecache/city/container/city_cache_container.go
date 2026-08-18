package container

import (
	"log"

	cachestore "cache-service/internal/cachestore"
	"cache-service/internal/simplecache/city/client"
	"cache-service/internal/simplecache/city/handler"
	serviceimpl "cache-service/internal/simplecache/city/service/impl"
)

type CityCacheContainer struct {
	CityCacheHandler *handler.CityCacheHandler
	IDToCityStore    *cachestore.Cache[int, string]
	SlugToIDStore    *cachestore.Cache[string, int]
}

func NewCityCacheContainer(cdnBaseURL string) *CityCacheContainer {
	idToCity := cachestore.New[int, string](
		cachestore.WithTTL(cachestore.NoExpiration),
	)
	slugToID := cachestore.New[string, int](
		cachestore.WithTTL(cachestore.NoExpiration),
	)

	cdnClient := client.NewCDNCityClient(cdnBaseURL)
	svc := serviceimpl.NewCityCacheService(idToCity, slugToID, cdnClient)
	h := handler.NewCityCacheHandler(svc)
	log.Printf("City cache configured (cdn=%s)", cdnBaseURL)

	return &CityCacheContainer{
		CityCacheHandler: h,
		IDToCityStore:    idToCity,
		SlugToIDStore:    slugToID,
	}
}
