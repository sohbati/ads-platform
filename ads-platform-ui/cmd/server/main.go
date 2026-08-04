package main

import (
	"log"

	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/container"
	"ads-platform-ui/internal/core/router"
)

type Application struct {
	config *config.Config
	router *router.Router
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Initialize() {
	app.config = config.Load()
	appContainer, err := container.NewAppContainer(app.config)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	app.router = router.NewRouter(appContainer)
	log.Println("Configuration and routes loaded")
}

func (app *Application) Run() error {
	ginRouter := app.router.SetupRoutes()
	log.Printf("Starting ads-platform-ui on port %s", app.config.Port)
	return ginRouter.Run(":" + app.config.Port)
}

func main() {
	app := NewApplication()
	app.Initialize()

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
