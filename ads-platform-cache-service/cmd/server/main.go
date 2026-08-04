package main

import (
	"log"

	"cache-service/internal/core/config"
	"cache-service/internal/core/container"
	"cache-service/internal/core/router"
)

// @title ads platform Service API
// @version 1.0
// @description Business Logic Service for ads platform Application
// @host localhost:8093
// @BasePath /

type Application struct {
	config    *config.Config
	router    *router.Router
	container *container.AppContainer
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Initialize() error {
	app.loadConfig()

	container := container.NewAppContainer()
	app.container = container
	app.router = router.NewRouter(container)
	log.Println("Router configured")

	return nil
}

func (app *Application) loadConfig() {
	app.config = config.Load()
	log.Println("Configuration loaded")
}

func (app *Application) Run() error {
	ginRouter := app.router.SetupRoutes()
	log.Printf("Starting ads platform Service on port %s", app.config.Port)
	return ginRouter.Run(":" + app.config.Port)
}

func main() {
	app := NewApplication()
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	defer app.container.OtpCacheContainer.OtpCacheStore.Stop()

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
