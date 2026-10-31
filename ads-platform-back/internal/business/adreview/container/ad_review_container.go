package container

import (
	"ads-platform/internal/business/adreview/handler"
	repoimpl "ads-platform/internal/business/adreview/repository/impl"
	serviceimpl "ads-platform/internal/business/adreview/service/impl"
	searchclient "ads-platform/internal/business/search/client"

	"gorm.io/gorm"
)

type AdReviewContainer struct {
	Handler *handler.AdReviewHandler
}

func NewAdReviewContainer(db *gorm.DB, cacheServiceURL string) *AdReviewContainer {
	svc := serviceimpl.NewAdReviewService(
		repoimpl.NewAdReviewRepository(db),
		searchclient.NewCatalogClient(cacheServiceURL, nil),
	)
	return &AdReviewContainer{Handler: handler.NewAdReviewHandler(svc)}
}
