package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	DatabaseType  string
	BrokerURL     string
	BrokerHTTPURL string
	StatsSubject  string
	StatsQueue    string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:          getEnv("PORT", "8099"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		DatabaseType:  getEnv("DATABASE_TYPE", "postgres"),
		BrokerURL:     os.Getenv("BROKER_URL"),
		BrokerHTTPURL: getEnv("BROKER_HTTP_URL", "http://localhost:8095"),
		StatsSubject:  getEnv("STATS_SUBJECT", "ads.stats.event"),
		StatsQueue:    getEnv("STATS_QUEUE", "ads-stats"),
	}

	brokerURL, err := resolveBrokerURL(cfg.BrokerURL, cfg.BrokerHTTPURL)
	if err != nil {
		return nil, err
	}
	cfg.BrokerURL = brokerURL
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
