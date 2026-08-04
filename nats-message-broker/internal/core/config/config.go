package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	NatsHost     string
	NatsPort     string
	NatsHTTPPort string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:         getEnv("PORT", "8095"),
		NatsHost:     getEnv("NATS_HOST", "127.0.0.1"),
		NatsPort:     getEnv("NATS_PORT", "-1"),
		NatsHTTPPort: getEnv("NATS_HTTP_PORT", "-1"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
