package main

import (
	"fmt"
	"log"

	"ads-platform/internal/core/config"
	"ads-platform/internal/core/container"
	"ads-platform/internal/core/db"
	"ads-platform/internal/core/router"

	"gorm.io/gorm"
)

// @title ads platform Service API
// @version 1.0
// @description Business Logic Service for ads platform Application
// @host localhost:8092
// @BasePath /

type Application struct {
	config   *config.Config
	database *gorm.DB
	router   *router.Router
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) Initialize() error {
	app.loadConfig()

	err := app.configDatabase()
	if err != nil {
		return err
	}
	container := container.NewAppContainer(app.database, app.config)

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
	log.Printf("Starting ads platform Service on port %s", app.config.ApplicationServerPort)
	return ginRouter.Run(":" + app.config.ApplicationServerPort)
}

func (app *Application) configDatabase() error {
	log.Println("Connecting to database...")

	database, err := db.Connect(app.config.DatabaseURL, app.config.DatabaseType)
	if err != nil {
		return fmt.Errorf("database unavailable: %w", err)
	}

	app.database = database
	return nil
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
