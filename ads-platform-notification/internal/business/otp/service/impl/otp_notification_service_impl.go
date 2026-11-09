package impl

import (
	"context"
	"fmt"
	"log"

	"ads-platform-notification/internal/business/otp/model"
	"ads-platform-notification/internal/business/otp/service"
	"ads-platform-notification/internal/core/sms"
)

type otpNotificationService struct {
	sender sms.Sender
}

func NewOtpNotificationService(sender sms.Sender) service.OtpNotificationService {
	return &otpNotificationService{sender: sender}
}

func (s *otpNotificationService) HandleOtpEvent(ctx context.Context, event *model.OtpEvent) error {
	if event.Mobile == "" {
		return fmt.Errorf("otp event: mobile is required")
	}
	if event.Otp == "" {
		return fmt.Errorf("otp event: otp is required")
	}

	res, err := s.sender.SendTemplate(ctx, sms.TemplateMessage{
		Mobile:   event.Mobile,
		Template: sms.TemplateOTP,
		Params:   map[string]string{sms.ParamCode: event.Otp},
	})
	if err != nil {
		return err
	}
	log.Printf("[OTP] sent mobile=%s provider=%s messageId=%s", event.Mobile, res.Provider, res.MessageID)
	return nil
}
