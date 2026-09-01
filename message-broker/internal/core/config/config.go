package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	BrokerHost        string
	BrokerPort        string
	BrokerMonitorPort string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:              getEnv("PORT", "8095"),
		BrokerHost:        getEnv("BROKER_HOST", "127.0.0.1"),
		BrokerPort:        getEnv("BROKER_PORT", "-1"),
		BrokerMonitorPort: getEnv("BROKER_MONITOR_PORT", "-1"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
