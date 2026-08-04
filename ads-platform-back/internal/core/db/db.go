package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const pingTimeout = 5 * time.Second

// Connect opens the database and verifies it is reachable before returning.
// Startup must treat any returned error as fatal.
func Connect(databaseURL, databaseType string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch databaseType {
	case "mysql":
		dialector = mysql.Open(databaseURL)
	case "postgres":
		dialector = postgres.Open(databaseURL)
	default:
		return nil, fmt.Errorf("unsupported database type %q (expected mysql or postgres)", databaseType)
	}

	database, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := VerifyConnection(database); err != nil {
		return nil, fmt.Errorf("verify database connection: %w", err)
	}

	log.Printf("Database connected and verified (%s)", databaseType)
	return database, nil
}

// VerifyConnection pings the database to ensure it is available.
func VerifyConnection(database *gorm.DB) error {
	sqlDB, err := database.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	return nil
}

func Migrate(db *gorm.DB) error {
	log.Println("Database migrations disabled - using manual SQL migrations")
	return nil
}
