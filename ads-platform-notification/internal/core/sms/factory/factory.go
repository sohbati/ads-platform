// Package factory builds the configured sms.Sender. Adding a provider means
// writing a package under sms/provider and adding a case in newProvider.
package factory

import (
	"fmt"
	"log"
	"strings"

	"ads-platform-notification/internal/core/config"
	"ads-platform-notification/internal/core/sms"
	"ads-platform-notification/internal/core/sms/allowlist"
	"ads-platform-notification/internal/core/sms/provider/logsms"
	"ads-platform-notification/internal/core/sms/provider/smsir"
)

// New builds the provider from SMS_PROVIDER. Outside production it is always
// wrapped in the recipient allowlist, so test traffic cannot reach customers.
func New(cfg *config.Config) (sms.Sender, error) {
	sender, err := newProvider(cfg.SMS)
	if err != nil {
		return nil, err
	}
	if cfg.IsProduction() {
		return sender, nil
	}
	if len(cfg.SMS.Allowlist) == 0 {
		log.Printf("WARNING: APP_ENV=%s and SMS_ALLOWLIST is empty, so every SMS will be blocked", cfg.AppEnv)
	}
	return allowlist.Wrap(sender, cfg.SMS.Allowlist)
}

func newProvider(cfg config.SMSConfig) (sms.Sender, error) {
	switch strings.ToLower(cfg.Provider) {
	case "", "log":
		return logsms.New(), nil
	case "smsir":
		return smsir.New(smsir.Config{
			APIKey:     cfg.SmsIR.APIKey,
			BaseURL:    cfg.SmsIR.BaseURL,
			LineNumber: cfg.SmsIR.LineNumber,
			Templates: map[string]smsir.Template{
				sms.TemplateOTP: {
					ID:     cfg.SmsIR.OtpTemplateID,
					Params: map[string]string{sms.ParamCode: cfg.SmsIR.OtpCodeParam},
				},
			},
		})
	default:
		return nil, fmt.Errorf("unknown SMS_PROVIDER %q (supported: log, smsir)", cfg.Provider)
	}
}
