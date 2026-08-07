package impl

import (
	"context"
	"log"

	"ads-platform-notification/internal/business/otp/sms"
)

type logProvider struct{}

func NewLogProvider() sms.Provider {
	return &logProvider{}
}

func (p *logProvider) SendSMS(ctx context.Context, mobile string, message string) error {
	log.Printf("[SMS] to=%s message=%q", mobile, message)
	return nil
}
