package main

import (
	"log"

	"ads-platform-crm/internal/core/config"
	"ads-platform-crm/internal/core/container"
	"ads-platform-crm/internal/core/router"
)

func main() {
	cfg := config.Load()
	log.Println("Configuration loaded")

	app := container.NewAppContainer(cfg)
	r := router.NewRouter(app)

	log.Printf("Starting ads-platform-crm on port %s", cfg.Port)
	if err := r.SetupRoutes().Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
