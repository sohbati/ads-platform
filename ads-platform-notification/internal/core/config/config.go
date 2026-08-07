package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	NatsURL       string
	NatsBrokerURL string
	OtpSubject    string
}

func Load() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		Port:          getEnv("PORT", "8096"),
		NatsURL:       os.Getenv("NATS_URL"),
		NatsBrokerURL: getEnv("NATS_BROKER_URL", "http://localhost:8095"),
		OtpSubject:    getEnv("OTP_SUBJECT", "notifications.otp.send"),
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
