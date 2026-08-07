package container

import (
	otpContainer "ads-platform-notification/internal/business/otp/container"
	"ads-platform-notification/internal/core/config"
	"ads-platform-notification/internal/core/natsconn"
)

type AppContainer struct {
	Nats *natsconn.Connection
	Otp  *otpContainer.OtpContainer
}

func NewAppContainer(cfg *config.Config, natsConn *natsconn.Connection) (*AppContainer, error) {
	otp, err := otpContainer.NewOtpContainer(natsConn, cfg.OtpSubject)
	if err != nil {
		return nil, err
	}

	return &AppContainer{
		Nats: natsConn,
		Otp:  otp,
	}, nil
}
