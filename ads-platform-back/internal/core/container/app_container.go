package container

import (
	otpContainer "ads-platform/internal/business/otp/container"
	searchContainer "ads-platform/internal/business/search/container"
	userContainer "ads-platform/internal/business/user/container"
	"ads-platform/internal/core/config"

	"gorm.io/gorm"
)

type AppContainer struct {
	User   *userContainer.UserContainer
	Otp    *otpContainer.OtpContainer
	Search *searchContainer.SearchContainer
}

func NewAppContainer(db *gorm.DB, cfg *config.Config) *AppContainer {
	return &AppContainer{
		User:   userContainer.NewUserContainer(db, cfg.DefaultCountryCode),
		Otp:    otpContainer.NewOtpContainer(cfg.CacheServiceURL, cfg.NatsURL, cfg.OtpSubject, cfg.DefaultCountryCode),
		Search: searchContainer.NewSearchContainer(db, cfg.CacheServiceURL),
	}
}
