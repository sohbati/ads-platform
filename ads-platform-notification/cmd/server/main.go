package main

import (
	"log"

	"ads-platform-notification/internal/core/config"
	"ads-platform-notification/internal/core/container"
	"ads-platform-notification/internal/core/natsconn"
	"ads-platform-notification/internal/core/router"
)

type Application struct {
	config    *config.Config
	nats      *natsconn.Connection
	container *container.AppContainer
	router    *router.Router
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Initialize() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	app.config = cfg
	log.Println("Configuration loaded")
	log.Printf("Using NATS at %s", app.config.NatsURL)

	natsConn, err := natsconn.Connect(app.config.NatsURL)
	if err != nil {
		return err
	}
	app.nats = natsConn

	appContainer, err := container.NewAppContainer(app.config, app.nats)
	if err != nil {
		natsConn.Close()
		return err
	}
	app.container = appContainer

	app.router = router.NewRouter(app.container)
	log.Println("Router configured")

	return nil
}

func (app *Application) Run() error {
	ginRouter := app.router.SetupRoutes()
	log.Printf("Starting ads-platform-notification on port %s", app.config.Port)
	return ginRouter.Run(":" + app.config.Port)
}

func main() {
	app := NewApplication()
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	defer app.container.Otp.OtpListener.Stop()
	defer app.nats.Close()

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
