package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	DatabaseType string
	NatsURL      string
	NatsBrokerURL string
	StatsSubject string
	StatsQueue   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:          getEnv("PORT", "8099"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		DatabaseType:  getEnv("DATABASE_TYPE", "postgres"),
		NatsURL:       os.Getenv("NATS_URL"),
		NatsBrokerURL: getEnv("NATS_BROKER_URL", "http://localhost:8095"),
		StatsSubject:  getEnv("STATS_SUBJECT", "ads.stats.event"),
		StatsQueue:    getEnv("STATS_QUEUE", "ads-stats"),
	}

	natsURL, err := resolveNatsURL(cfg.NatsURL, cfg.NatsBrokerURL)
	if err != nil {
		return nil, err
	}
	cfg.NatsURL = natsURL
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
