package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string
	AppName string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:    getEnv("PORT", "8100"),
		AppName: getEnv("APP_NAME", "ruab CRM"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
