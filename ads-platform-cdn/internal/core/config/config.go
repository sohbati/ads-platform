package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	StaticDir      string
	CategoryJSON   string
	CitiesJSON     string
	LocationJSON   string
}

func Load() *Config {
	_ = godotenv.Load()

	staticDir := getEnv("STATIC_DIR", "./cdn")

	return &Config{
		Port:         getEnv("PORT", "4000"),
		StaticDir:    staticDir,
		CategoryJSON: getEnv("CATEGORY_JSON", filepath.Join(staticDir, "json", "category.json")),
		CitiesJSON:   getEnv("CITIES_JSON", filepath.Join(staticDir, "json", "cities.json")),
		LocationJSON: getEnv("LOCATION_JSON", filepath.Join(staticDir, "json", "location.json")),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
