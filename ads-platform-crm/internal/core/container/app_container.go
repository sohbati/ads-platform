package container

import "ads-platform-crm/internal/core/config"

type AppContainer struct {
	Config *config.Config
}

func NewAppContainer(cfg *config.Config) *AppContainer {
	return &AppContainer{Config: cfg}
}
