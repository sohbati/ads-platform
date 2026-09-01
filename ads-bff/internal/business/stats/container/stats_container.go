package container

import (
	"log"
	"time"

	"ads-bff/internal/business/auth/service"
	"ads-bff/internal/business/stats/handler"
	"ads-bff/internal/business/stats/publisher"
	"ads-bff/internal/core/config"
	"ads-bff/internal/core/ratelimit"
)

type StatsContainer struct {
	Handler *handler.StatsHandler
}

func NewStatsContainer(cfg *config.Config, auth service.AuthService) *StatsContainer {
	pub, err := publisher.New(cfg.NatsURL, cfg.StatsSubject)
	if err != nil {
		log.Printf("stats NATS publisher disabled: %v", err)
		pub = publisher.Noop()
	}
	limiter := ratelimit.New(time.Minute, cfg.StatsRatePerMin)
	return &StatsContainer{
		Handler: handler.NewStatsHandler(cfg, auth, pub, limiter),
	}
}
