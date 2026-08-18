package container

import (
	"log"

	cachestore "cache-service/internal/cachestore"
	"cache-service/internal/simplecache/category/client"
	"cache-service/internal/simplecache/category/handler"
	serviceimpl "cache-service/internal/simplecache/category/service/impl"
)

type CategoryCacheContainer struct {
	CategoryCacheHandler *handler.CategoryCacheHandler
	SlugToCategoryStore  *cachestore.Cache[string, string]
	IDToSlugStore        *cachestore.Cache[int, string]
}

func NewCategoryCacheContainer(cdnBaseURL string) *CategoryCacheContainer {
	slugToCategory := cachestore.New[string, string](
		cachestore.WithTTL(cachestore.NoExpiration),
	)
	idToSlug := cachestore.New[int, string](
		cachestore.WithTTL(cachestore.NoExpiration),
	)

	cdnClient := client.NewCDNCategoryClient(cdnBaseURL)
	svc := serviceimpl.NewCategoryCacheService(slugToCategory, idToSlug, cdnClient)
	h := handler.NewCategoryCacheHandler(svc)
	log.Printf("Category cache configured (cdn=%s)", cdnBaseURL)

	return &CategoryCacheContainer{
		CategoryCacheHandler: h,
		SlugToCategoryStore:  slugToCategory,
		IDToSlugStore:        idToSlug,
	}
}
