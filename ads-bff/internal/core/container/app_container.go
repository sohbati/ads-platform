package container

import (
	"fmt"

	authContainer "ads-bff/internal/business/auth/container"
	"ads-bff/internal/core/config"
	"ads-bff/internal/core/proxy"
)

type AppContainer struct {
	Config  *config.Config
	Backend *proxy.BackendProxy
	Auth    *authContainer.AuthContainer
}

func NewAppContainer(cfg *config.Config) (*AppContainer, error) {
	backend, err := proxy.NewBackendProxy(cfg.BackendAPIBaseURL)
	if err != nil {
		return nil, fmt.Errorf("init backend proxy: %w", err)
	}

	return &AppContainer{
		Config:  cfg,
		Backend: backend,
		Auth:    authContainer.NewAuthContainer(cfg),
	}, nil
}
