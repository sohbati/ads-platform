package container

import (
	"log"

	"ads-platform/internal/business/otp/client"
	"ads-platform/internal/business/otp/handler"
	serviceimpl "ads-platform/internal/business/otp/service/impl"
)

type OtpContainer struct {
	OtpHandler *handler.OtpHandler
}

func NewOtpContainer(cacheServiceURL string, natsURL string, otpSubject string) *OtpContainer {
	cacheClient := client.NewOtpCacheClient(cacheServiceURL, nil)

	eventPublisher, err := client.NewOtpEventPublisher(natsURL, otpSubject)
	if err != nil {
		log.Printf("OTP NATS publisher disabled: %v", err)
		eventPublisher, _ = client.NewOtpEventPublisher("", otpSubject)
	}

	otpService := serviceimpl.NewOtpService(cacheClient, eventPublisher)
	otpHandler := handler.NewOtpHandler(otpService)

	return &OtpContainer{
		OtpHandler: otpHandler,
	}
}
