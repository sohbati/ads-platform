package container

import (
	"ads-bff/internal/business/auth/service"
	"ads-bff/internal/business/profile/handler"
	backendclient "ads-bff/internal/core/client/backend"
	"ads-bff/internal/core/config"
)

type ProfileContainer struct {
	Handler *handler.ProfileHandler
}

func NewProfileContainer(cfg *config.Config, auth service.AuthService) *ProfileContainer {
	backend := backendclient.NewClient(cfg.BackendAPIBaseURL, nil)
	return &ProfileContainer{
		Handler: handler.NewProfileHandler(cfg, auth, backend),
	}
}
