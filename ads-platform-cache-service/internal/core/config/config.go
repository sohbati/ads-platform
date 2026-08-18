package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	fooEnv       string
	DatabaseType string
	Port         string
	CDNBaseURL   string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:       getEnv("PORT", "8093"),
		CDNBaseURL: getEnv("CDN_BASE_URL", "http://localhost:4000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
