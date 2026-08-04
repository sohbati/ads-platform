package container

import (
	"context"
	"fmt"
	"time"

	categoryContainer "ads-platform-ui/internal/business/category/container"
	locationContainer "ads-platform-ui/internal/business/location/container"
	myinfoContainer "ads-platform-ui/internal/business/myinfo/container"
	newadContainer "ads-platform-ui/internal/business/newad/container"
	queryadsContainer "ads-platform-ui/internal/business/queryads/container"
	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
)

type AppContainer struct {
	Config   *config.Config
	I18n     *i18n.Registry
	CDN      *cdn.Client
	Cities   *cities.Catalog
	QueryAds *queryadsContainer.QueryAdsContainer
	MyInfo   *myinfoContainer.MyInfoContainer
	NewAd    *newadContainer.NewAdContainer
	Category *categoryContainer.CategoryContainer
	Location *locationContainer.LocationContainer
}

func NewAppContainer(cfg *config.Config) (*AppContainer, error) {
	reg, err := i18n.LoadRegistry(i18n.DefaultLocale)
	if err != nil {
		return nil, err
	}

	cdnClient := cdn.NewClient(cfg.CDNBaseURL, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	catalog, err := cities.FetchFromCDN(ctx, cfg.CDNBaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("load cities from CDN: %w", err)
	}

	return &AppContainer{
		Config:   cfg,
		I18n:     reg,
		CDN:      cdnClient,
		Cities:   catalog,
		QueryAds: queryadsContainer.NewQueryAdsContainer(cfg, reg, catalog),
		MyInfo:   myinfoContainer.NewMyInfoContainer(cfg, reg, catalog),
		NewAd:    newadContainer.NewNewAdContainer(cfg, reg, catalog),
		Category: categoryContainer.NewCategoryContainer(cfg, reg, catalog, cdnClient),
		Location: locationContainer.NewLocationContainer(cfg, reg, catalog, cdnClient),
	}, nil
}
