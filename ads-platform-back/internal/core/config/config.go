package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL           string
	DatabaseType          string
	ApplicationServerPort string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		DatabaseType:          os.Getenv("DATABASE_TYPE"),
		ApplicationServerPort: os.Getenv("APPLICATION_SERVER_PORT"),
	}
}
