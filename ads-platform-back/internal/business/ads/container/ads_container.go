package container

import (
	"log"

	"ads-platform/internal/business/ads/handler"
	repoimpl "ads-platform/internal/business/ads/repository/impl"
	serviceimpl "ads-platform/internal/business/ads/service/impl"
	searchclient "ads-platform/internal/business/search/client"
	"ads-platform/internal/core/config"
	"ads-platform/internal/core/storage"

	"gorm.io/gorm"
)

type AdsContainer struct {
	AdHandler *handler.AdHandler
}

func NewAdsContainer(db *gorm.DB, cfg *config.Config) *AdsContainer {
	objects, err := storage.NewMinio(
		cfg.MinioEndpoint,
		cfg.MinioAccessKey,
		cfg.MinioSecretKey,
		cfg.MinioBucket,
		cfg.MinioPublicURL,
		cfg.MinioUseSSL,
	)
	if err != nil {
		log.Printf("minio unavailable (ads can still be created without pictures): %v", err)
		objects = nil
	}

	catalog := searchclient.NewCatalogClient(cfg.CacheServiceURL, nil)
	adRepo := repoimpl.NewAdRepository(db)
	imageRepo := repoimpl.NewAdImageRepository(db)
	svc := serviceimpl.NewAdService(adRepo, imageRepo, catalog, objects, cfg.MaxAdPictures, cfg.MaxAdPictureBytes)

	return &AdsContainer{
		AdHandler: handler.NewAdHandler(svc),
	}
}
