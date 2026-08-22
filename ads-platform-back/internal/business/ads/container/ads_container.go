package container

import (
	"ads-platform/internal/business/ads/handler"
	repoimpl "ads-platform/internal/business/ads/repository/impl"
	serviceimpl "ads-platform/internal/business/ads/service/impl"

	"gorm.io/gorm"
)

type AdsContainer struct {
	AdImageHandler *handler.AdImageHandler
}

func NewAdsContainer(db *gorm.DB) *AdsContainer {
	imageRepo := repoimpl.NewAdImageRepository(db)
	imageService := serviceimpl.NewAdImageService(imageRepo)

	return &AdsContainer{
		AdImageHandler: handler.NewAdImageHandler(imageService),
	}
}
