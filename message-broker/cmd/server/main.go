package main

import (
	"fmt"
	"log"
	"strconv"

	"message-broker/internal/core/broker"
	"message-broker/internal/core/brokerserver"
	"message-broker/internal/core/config"
	"message-broker/internal/core/container"
	"message-broker/internal/core/router"
)

// @title Message Broker API
// @version 1.0
// @description Event message broker service
// @host localhost:8095
// @BasePath /

type Application struct {
	config       *config.Config
	brokerServer *brokerserver.EmbeddedServer
	broker       *broker.Connection
	container    *container.AppContainer
	router       *router.Router
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Initialize() error {
	app.loadConfig()

	brokerPort, err := strconv.Atoi(app.config.BrokerPort)
	if err != nil {
		return fmt.Errorf("invalid BROKER_PORT %q: %w", app.config.BrokerPort, err)
	}

	monitorPort, err := strconv.Atoi(app.config.BrokerMonitorPort)
	if err != nil {
		return fmt.Errorf("invalid BROKER_MONITOR_PORT %q: %w", app.config.BrokerMonitorPort, err)
	}

	brokerServer, err := brokerserver.Start(app.config.BrokerHost, brokerPort, monitorPort)
	if err != nil {
		return err
	}
	app.brokerServer = brokerServer

	brokerConn, err := broker.Connect(brokerServer.ClientURL())
	if err != nil {
		brokerServer.Shutdown()
		return err
	}
	app.broker = brokerConn

	app.container = container.NewAppContainer(app.brokerServer, app.broker)
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
	log.Printf("Starting message-broker on port %s", app.config.Port)
	log.Printf("Broker monitoring panel: http://localhost:%s%s/", app.config.Port, brokerserver.MonitorPath)
	return ginRouter.Run(":" + app.config.Port)
}

func main() {
	app := NewApplication()
	if err := app.Initialize(); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	defer app.broker.Close()
	defer app.brokerServer.Shutdown()

	if err := app.Run(); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
