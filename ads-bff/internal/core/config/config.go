package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	BackendAPIBaseURL  string
	CacheServiceURL    string
	SessionCookieName  string
	SessionTTL         time.Duration
	CookieSecure       bool
	DefaultCountryCode string
	BrokerURL          string
	BrokerHTTPURL      string
	StatsSubject       string
	StatsRatePerMin    int
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:               getEnv("PORT", "8097"),
		BackendAPIBaseURL:  getEnv("BACKEND_API_BASE_URL", "http://localhost:8092"),
		CacheServiceURL:    getEnv("CACHE_SERVICE_URL", "http://localhost:8093"),
		SessionCookieName:  getEnv("SESSION_COOKIE_NAME", "ads_session"),
		SessionTTL:         getDurationEnv("SESSION_TTL", 24*time.Hour),
		CookieSecure:       getBoolEnv("COOKIE_SECURE", false),
		DefaultCountryCode: getEnv("DEFAULT_COUNTRY_CODE", "+98"),
		BrokerURL:          os.Getenv("BROKER_URL"),
		BrokerHTTPURL:      getEnv("BROKER_HTTP_URL", "http://localhost:8095"),
		StatsSubject:       getEnv("STATS_SUBJECT", "ads.stats.event"),
		StatsRatePerMin:    getIntEnv("STATS_RATE_PER_MIN", 60),
	}

	if brokerURL, err := resolveBrokerURL(cfg.BrokerURL, cfg.BrokerHTTPURL); err == nil && brokerURL != "" {
		cfg.BrokerURL = brokerURL
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		parsed, err := strconv.ParseBool(value)
		if err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if n, err := strconv.Atoi(value); err == nil && n > 0 {
			return n
		}
	}
	return defaultValue
}
