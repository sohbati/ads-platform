package container

import (
	"fmt"

	adsContainer "ads-bff/internal/business/ads/container"
	authContainer "ads-bff/internal/business/auth/container"
	profileContainer "ads-bff/internal/business/profile/container"
	statsContainer "ads-bff/internal/business/stats/container"
	"ads-bff/internal/core/config"
	"ads-bff/internal/core/proxy"
)

type AppContainer struct {
	Config  *config.Config
	Backend *proxy.BackendProxy
	Auth    *authContainer.AuthContainer
	Ads     *adsContainer.AdsContainer
	Profile *profileContainer.ProfileContainer
	Stats   *statsContainer.StatsContainer
}

func NewAppContainer(cfg *config.Config) (*AppContainer, error) {
	backend, err := proxy.NewBackendProxy(cfg.BackendAPIBaseURL)
	if err != nil {
		return nil, fmt.Errorf("init backend proxy: %w", err)
	}

	auth := authContainer.NewAuthContainer(cfg)
	return &AppContainer{
		Config:  cfg,
		Backend: backend,
		Auth:    auth,
		Ads:     adsContainer.NewAdsContainer(cfg, auth.AuthService),
		Profile: profileContainer.NewProfileContainer(cfg, auth.AuthService),
		Stats:   statsContainer.NewStatsContainer(cfg, auth.AuthService),
	}, nil
}
