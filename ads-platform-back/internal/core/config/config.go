package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL           string
	DatabaseType          string
	ApplicationServerPort string
	CacheServiceURL       string
	NatsURL               string
	NatsBrokerURL         string
	OtpSubject            string
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		DatabaseType:          os.Getenv("DATABASE_TYPE"),
		ApplicationServerPort: os.Getenv("APPLICATION_SERVER_PORT"),
		CacheServiceURL:       os.Getenv("CACHE_SERVICE_URL"),
		NatsURL:               os.Getenv("NATS_URL"),
		NatsBrokerURL:         getEnv("NATS_BROKER_URL", "http://localhost:8095"),
		OtpSubject:            getEnv("OTP_SUBJECT", "notifications.otp.send"),
	}

	natsURL, err := resolveNatsURL(cfg.NatsURL, cfg.NatsBrokerURL)
	if err != nil {
		log.Printf("NATS URL discovery failed: %v", err)
	} else {
		cfg.NatsURL = natsURL
		if natsURL != "" {
			log.Printf("Using NATS at %s", natsURL)
		}
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
