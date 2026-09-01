package main

import (
	"log"

	"ads-platform-stats/internal/core/config"
	"ads-platform-stats/internal/core/container"
	"ads-platform-stats/internal/core/db"
	"ads-platform-stats/internal/core/natsconn"
	"ads-platform-stats/internal/core/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println("Configuration loaded")

	database, err := db.Connect(cfg.DatabaseURL, cfg.DatabaseType)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	natsConn, err := natsconn.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatalf("Failed to connect NATS: %v", err)
	}
	defer natsConn.Close()

	stats, err := container.NewStatsContainer(cfg, natsConn, database)
	if err != nil {
		log.Fatalf("Failed to start stats listener: %v", err)
	}
	defer stats.Listener.Stop()

	r := router.NewRouter(stats, natsConn.IsConnected)
	log.Printf("Starting ads-platform-stats on port %s", cfg.Port)
	if err := r.SetupRoutes().Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
