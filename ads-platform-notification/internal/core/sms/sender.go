// Package sms is the provider-neutral SMS API used by the business code.
//
// Business code never talks to a vendor directly. It sends a message that names a
// logical template (for example TemplateOTP) and passes named parameters. Each
// provider maps those names to its own template IDs and parameter names, so you
// switch vendors by changing configuration, not code.
package sms

import (
	"context"
	"errors"
)

// Logical template names. Each provider maps these to its own template IDs.
const (
	TemplateOTP = "otp" // params: ParamCode
)

// Logical parameter names used in TemplateMessage.Params.
const (
	ParamCode = "code"
)

// TemplateMessage is a message sent through a template the provider has pre-approved
// (sms.ir "verify", Kavenegar "lookup", and so on). Iranian carriers need this for
// OTP messages to reach numbers that block advertising SMS.
type TemplateMessage struct {
	Mobile   string
	Template string
	Params   map[string]string
}

// TextMessage is a free-text message sent from the provider's sender line.
type TextMessage struct {
	Mobile string
	Text   string
}

// Sender is what every SMS provider implements.
type Sender interface {
	// Name identifies the provider in logs.
	Name() string
	SendTemplate(ctx context.Context, msg TemplateMessage) (*Result, error)
	SendText(ctx context.Context, msg TextMessage) (*Result, error)
}

// Result holds what a provider reports about an accepted message.
type Result struct {
	Provider  string
	MessageID string
}

var (
	// ErrUnknownTemplate means the provider has no mapping for the logical template.
	ErrUnknownTemplate = errors.New("sms: template not configured for provider")
	// ErrNotSupported means the provider cannot send this kind of message with the current configuration.
	ErrNotSupported = errors.New("sms: operation not supported by provider")
	// ErrRecipientNotAllowed means the allowlist blocked the number (non-production environments).
	ErrRecipientNotAllowed = errors.New("sms: recipient not in allowlist")
)

// ProviderError is a failure the provider reported.
type ProviderError struct {
	Provider   string
	HTTPStatus int
	Code       int
	Message    string
	// Retryable is true for transient failures (rate limits, 5xx, network errors).
	Retryable bool
}

func (e *ProviderError) Error() string {
	return "sms " + e.Provider + ": " + e.Message
}
