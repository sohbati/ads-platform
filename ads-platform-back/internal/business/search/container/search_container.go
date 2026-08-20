package container

import (
	"ads-platform/internal/business/search/client"
	"ads-platform/internal/business/search/handler"
	repoimpl "ads-platform/internal/business/search/repository/impl"
	serviceimpl "ads-platform/internal/business/search/service/impl"

	"gorm.io/gorm"
)

type SearchContainer struct {
	SearchHandler *handler.SearchHandler
}

func NewSearchContainer(db *gorm.DB, cacheServiceURL string) *SearchContainer {
	catalog := client.NewCatalogClient(cacheServiceURL, nil)
	repo := repoimpl.NewSearchRepository(db)
	svc := serviceimpl.NewSearchService(repo, catalog)
	h := handler.NewSearchHandler(svc)
	return &SearchContainer{SearchHandler: h}
}
