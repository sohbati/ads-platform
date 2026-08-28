package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	AppName            string
	DefaultCity        string
	CDNBaseURL         string
	MediaCDNURL        string
	BFFBaseURL         string
	DefaultCountryCode string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8094"),
		AppName:     getEnv("APP_NAME", "Ruab"),
		DefaultCity: getEnv("DEFAULT_CITY", "tehran"),
		CDNBaseURL:         getEnv("CDN_BASE_URL", "http://localhost:4000"),
		MediaCDNURL:        getEnv("MEDIA_CDN_URL", "http://localhost:8098"),
		BFFBaseURL:         getEnv("BFF_BASE_URL", "http://localhost:8097"),
		DefaultCountryCode: getEnv("DEFAULT_COUNTRY_CODE", "+98"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
