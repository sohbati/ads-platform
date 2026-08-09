package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	BackendAPIBaseURL string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:              getEnv("PORT", "8097"),
		BackendAPIBaseURL: getEnv("BACKEND_API_BASE_URL", "http://localhost:8092"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
