package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	AppName           string
	DefaultCity       string
	CDNBaseURL        string
	BackendAPIBaseURL string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:              getEnv("PORT", "8094"),
		AppName:           getEnv("APP_NAME", "Ruab"),
		DefaultCity:       getEnv("DEFAULT_CITY", "tehran"),
		CDNBaseURL:        getEnv("CDN_BASE_URL", "http://localhost:4000"),
		BackendAPIBaseURL: getEnv("BACKEND_API_BASE_URL", "http://localhost:8092"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
