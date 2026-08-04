package service

import (
	"context"

	"ads-platform/internal/business/otp/model"
)

type OtpService interface {
	SendOTP(ctx context.Context, mobile string) (*model.SendOtpResponse, error)
	VerifyOTP(ctx context.Context, mobile string, otp string) (*model.VerifyOtpResponse, error)
}
