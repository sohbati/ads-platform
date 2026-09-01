package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	DefaultCountryCode    string
	OtpResendAfter        time.Duration
	MaxAdPictures         int
	MaxAdPictureBytes     int64
	MinioEndpoint         string
	MinioAccessKey        string
	MinioSecretKey        string
	MinioBucket           string
	MinioUseSSL           bool
	MinioPublicURL        string
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		DatabaseType:          os.Getenv("DATABASE_TYPE"),
		ApplicationServerPort: os.Getenv("APPLICATION_SERVER_PORT"),
		CacheServiceURL:       os.Getenv("CACHE_SERVICE_URL"),
		NatsURL:               os.Getenv("NATS_URL"),
		NatsBrokerURL:         os.Getenv("NATS_BROKER_URL"),
		OtpSubject:            os.Getenv("OTP_SUBJECT"),
		DefaultCountryCode:    os.Getenv("DEFAULT_COUNTRY_CODE"),
		OtpResendAfter:        envDurationSeconds("OTP_RESEND_SECONDS", 60),
		MaxAdPictures:         envInt("ADS_MAX_PICTURES"),
		MaxAdPictureBytes:     int64(envInt("ADS_MAX_PICTURE_BYTES")),
		MinioEndpoint:         os.Getenv("MINIO_ENDPOINT"),
		MinioAccessKey:        os.Getenv("MINIO_ACCESS_KEY"),
		MinioSecretKey:        os.Getenv("MINIO_SECRET_KEY"),
		MinioBucket:           os.Getenv("MINIO_BUCKET"),
		MinioUseSSL:           os.Getenv("MINIO_USE_SSL") == "true",
		MinioPublicURL:        os.Getenv("MINIO_PUBLIC_URL"),
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

func envInt(key string) int {
	v, _ := strconv.Atoi(os.Getenv(key))
	return v
}

func envDurationSeconds(key string, defaultSeconds int) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return time.Duration(defaultSeconds) * time.Second
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return time.Duration(defaultSeconds) * time.Second
	}
	return time.Duration(n) * time.Second
}
