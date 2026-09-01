package container

import (
	"ads-platform-stats/internal/business/stats/listener"
	repoimpl "ads-platform-stats/internal/business/stats/repository/impl"
	"ads-platform-stats/internal/business/stats/service"
	"ads-platform-stats/internal/core/config"
	"ads-platform-stats/internal/core/natsconn"

	"gorm.io/gorm"
)

type StatsContainer struct {
	Listener *listener.StatsListener
}

func NewStatsContainer(cfg *config.Config, natsConn *natsconn.Connection, db *gorm.DB) (*StatsContainer, error) {
	svc := service.NewStatsService(repoimpl.NewStatsRepository(db))
	lis, err := listener.NewStatsListener(natsConn, cfg.StatsSubject, cfg.StatsQueue, svc)
	if err != nil {
		return nil, err
	}
	return &StatsContainer{Listener: lis}, nil
}
