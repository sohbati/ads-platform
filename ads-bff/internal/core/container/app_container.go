package container

import (
	"fmt"

	"ads-bff/internal/core/config"
	"ads-bff/internal/core/proxy"
)

type AppContainer struct {
	Config  *config.Config
	Backend *proxy.BackendProxy
}

func NewAppContainer(cfg *config.Config) (*AppContainer, error) {
	backend, err := proxy.NewBackendProxy(cfg.BackendAPIBaseURL)
	if err != nil {
		return nil, fmt.Errorf("init backend proxy: %w", err)
	}

	return &AppContainer{
		Config:  cfg,
		Backend: backend,
	}, nil
}
