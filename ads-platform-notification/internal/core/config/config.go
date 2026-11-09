package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Environments. Anything other than production counts as a test environment.
const (
	EnvLocal      = "local"
	EnvUAT        = "uat"
	EnvProduction = "production"
)

type Config struct {
	AppEnv        string // APP_ENV: local (default) | uat | production
	Port          string
	BrokerURL     string
	BrokerHTTPURL string
	OtpSubject    string
	SMS           SMSConfig
}

// SMSConfig selects the SMS provider and holds the settings of each supported one.
type SMSConfig struct {
	Provider string // SMS_PROVIDER: log (default) | smsir
	// Allowlist (SMS_ALLOWLIST, comma separated) is always enforced outside
	// production. Empty means nobody receives SMS. Ignored in production.
	Allowlist []string
	SmsIR     SmsIRConfig
}

type SmsIRConfig struct {
	APIKey        string
	BaseURL       string
	LineNumber    int64
	OtpTemplateID int
	OtpCodeParam  string
}

func Load() (*Config, error) {
	godotenv.Load()

	appEnv := strings.ToLower(getEnv("APP_ENV", EnvLocal))
	switch appEnv {
	case EnvLocal, EnvUAT, EnvProduction:
	default:
		return nil, fmt.Errorf("APP_ENV %q is invalid (use local, uat or production)", appEnv)
	}

	cfg := &Config{
		AppEnv:        appEnv,
		Port:          getEnv("PORT", "8096"),
		BrokerURL:     os.Getenv("BROKER_URL"),
		BrokerHTTPURL: getEnv("BROKER_HTTP_URL", "http://localhost:8095"),
		OtpSubject:    getEnv("OTP_SUBJECT", "notifications.otp.send"),
	}

	sms, err := loadSMSConfig()
	if err != nil {
		return nil, err
	}
	cfg.SMS = sms

	brokerURL, err := resolveBrokerURL(cfg.BrokerURL, cfg.BrokerHTTPURL)
	if err != nil {
		return nil, err
	}
	cfg.BrokerURL = brokerURL

	return cfg, nil
}

// IsProduction is true only when APP_ENV=production is set explicitly.
func (c *Config) IsProduction() bool { return c.AppEnv == EnvProduction }

func loadSMSConfig() (SMSConfig, error) {
	templateID, err := strconv.Atoi(getEnv("SMSIR_OTP_TEMPLATE_ID", "123456"))
	if err != nil {
		return SMSConfig{}, fmt.Errorf("SMSIR_OTP_TEMPLATE_ID: %w", err)
	}
	var line int64
	if v := os.Getenv("SMSIR_LINE_NUMBER"); v != "" {
		if line, err = strconv.ParseInt(v, 10, 64); err != nil {
			return SMSConfig{}, fmt.Errorf("SMSIR_LINE_NUMBER: %w", err)
		}
	}
	return SMSConfig{
		Provider:  getEnv("SMS_PROVIDER", "log"),
		Allowlist: splitList(os.Getenv("SMS_ALLOWLIST")),
		SmsIR: SmsIRConfig{
			APIKey:        os.Getenv("SMSIR_API_KEY"),
			BaseURL:       getEnv("SMSIR_BASE_URL", "https://api.sms.ir"),
			LineNumber:    line,
			OtpTemplateID: templateID,
			OtpCodeParam:  getEnv("SMSIR_OTP_CODE_PARAM", "Code"),
		},
	}, nil
}

func splitList(v string) []string {
	var out []string
	for _, item := range strings.Split(v, ",") {
		if item = strings.TrimSpace(item); item != "" {
			out = append(out, item)
		}
	}
	return out
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
