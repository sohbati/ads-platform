package main

import (
	"log"

	"ads-bff/internal/core/config"
	"ads-bff/internal/core/container"
	"ads-bff/internal/core/router"
)

type Application struct {
	config    *config.Config
	router    *router.Router
	container *container.AppContainer
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Initialize() error {
	app.config = config.Load()

	c, err := container.NewAppContainer(app.config)
	if err != nil {
		return err
	}
	app.container = c
	app.router = router.NewRouter(c)

	log.Println("Router configured")
	return nil
}

func (app *Application) Run() error {
	ginRouter := app.router.SetupRoutes()
	log.Printf("Starting ads-bff on port %s (backend: %s)", app.config.Port, app.config.BackendAPIBaseURL)
	return ginRouter.Run(":" + app.config.Port)
}

func main() {
	app := NewApplication()
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
