package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	BrokerURL     string
	BrokerHTTPURL string
	OtpSubject    string
}

func Load() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		Port:          getEnv("PORT", "8096"),
		BrokerURL:     os.Getenv("BROKER_URL"),
		BrokerHTTPURL: getEnv("BROKER_HTTP_URL", "http://localhost:8095"),
		OtpSubject:    getEnv("OTP_SUBJECT", "notifications.otp.send"),
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
