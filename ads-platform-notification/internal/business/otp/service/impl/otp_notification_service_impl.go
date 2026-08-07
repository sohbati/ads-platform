package impl

import (
	"context"
	"fmt"
	"log"

	"ads-platform-notification/internal/business/otp/model"
	"ads-platform-notification/internal/business/otp/service"
	"ads-platform-notification/internal/business/otp/sms"
)

type otpNotificationService struct {
	smsProvider sms.Provider
}

func NewOtpNotificationService(smsProvider sms.Provider) service.OtpNotificationService {
	return &otpNotificationService{smsProvider: smsProvider}
}

func (s *otpNotificationService) HandleOtpEvent(ctx context.Context, event *model.OtpEvent) error {
	if event.Mobile == "" {
		return fmt.Errorf("otp event: mobile is required")
	}
	if event.Otp == "" {
		return fmt.Errorf("otp event: otp is required")
	}

	log.Printf("[OTP] received mobile=%s otp=%s", event.Mobile, event.Otp)

	message := fmt.Sprintf("Your verification code is: %s", event.Otp)
	return s.smsProvider.SendSMS(ctx, event.Mobile, message)
}
