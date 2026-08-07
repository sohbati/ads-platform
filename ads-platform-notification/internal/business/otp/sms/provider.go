package sms

import "context"

type Provider interface {
	SendSMS(ctx context.Context, mobile string, message string) error
}
