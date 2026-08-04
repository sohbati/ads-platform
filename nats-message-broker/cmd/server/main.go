package main

import (
	"fmt"
	"log"
	"strconv"

	"nats-message-broker/internal/core/config"
	"nats-message-broker/internal/core/container"
	"nats-message-broker/internal/core/natsconn"
	"nats-message-broker/internal/core/natsserver"
	"nats-message-broker/internal/core/router"
)

// @title NATS Message Broker API
// @version 1.0
// @description Event message broker service with embedded NATS server
// @host localhost:8095
// @BasePath /

type Application struct {
	config     *config.Config
	natsServer *natsserver.EmbeddedServer
	nats       *natsconn.Connection
	container  *container.AppContainer
	router     *router.Router
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Initialize() error {
	app.loadConfig()

	natsPort, err := strconv.Atoi(app.config.NatsPort)
	if err != nil {
		return fmt.Errorf("invalid NATS_PORT %q: %w", app.config.NatsPort, err)
	}

	natsHTTPPort, err := strconv.Atoi(app.config.NatsHTTPPort)
	if err != nil {
		return fmt.Errorf("invalid NATS_HTTP_PORT %q: %w", app.config.NatsHTTPPort, err)
	}

	natsServer, err := natsserver.Start(app.config.NatsHost, natsPort, natsHTTPPort)
	if err != nil {
		return err
	}
	app.natsServer = natsServer

	natsConn, err := natsconn.Connect(natsServer.ClientURL())
	if err != nil {
		natsServer.Shutdown()
		return err
	}
	app.nats = natsConn

	app.container = container.NewAppContainer(app.natsServer, app.nats)
	app.router = router.NewRouter(app.container)
	log.Println("Router configured")

	return nil
}

func (app *Application) loadConfig() {
	app.config = config.Load()
	log.Println("Configuration loaded")
}

func (app *Application) Run() error {
	ginRouter := app.router.SetupRoutes()
	log.Printf("Starting nats-message-broker on port %s", app.config.Port)
	log.Printf("NATS monitoring panel: http://localhost:%s%s/", app.config.Port, natsserver.MonitorPath)
	return ginRouter.Run(":" + app.config.Port)
}

func main() {
	app := NewApplication()
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	defer app.nats.Close()
	defer app.natsServer.Shutdown()

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
