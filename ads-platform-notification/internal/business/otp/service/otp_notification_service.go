package service

import (
	"context"

	"ads-platform-notification/internal/business/otp/model"
)

type OtpNotificationService interface {
	HandleOtpEvent(ctx context.Context, event *model.OtpEvent) error
}
