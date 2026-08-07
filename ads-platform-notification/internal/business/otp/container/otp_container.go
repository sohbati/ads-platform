package container

import (
	"ads-platform-notification/internal/business/otp/listener"
	serviceimpl "ads-platform-notification/internal/business/otp/service/impl"
	smsimpl "ads-platform-notification/internal/business/otp/sms/impl"
	"ads-platform-notification/internal/core/natsconn"
)

type OtpContainer struct {
	OtpListener *listener.OtpListener
}

func NewOtpContainer(natsConn *natsconn.Connection, subject string) (*OtpContainer, error) {
	smsProvider := smsimpl.NewLogProvider()
	otpService := serviceimpl.NewOtpNotificationService(smsProvider)

	otpListener, err := listener.NewOtpListener(natsConn, subject, otpService)
	if err != nil {
		return nil, err
	}

	return &OtpContainer{
		OtpListener: otpListener,
	}, nil
}
