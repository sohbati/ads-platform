package container

import (
	"log"

	"ads-platform/internal/business/otp/client"
	"ads-platform/internal/business/otp/handler"
	serviceimpl "ads-platform/internal/business/otp/service/impl"
	"ads-platform/internal/core/mobile"
)

type OtpContainer struct {
	OtpHandler *handler.OtpHandler
}

func NewOtpContainer(cacheServiceURL, natsURL, otpSubject, defaultCountryCode string) *OtpContainer {
	cacheClient := client.NewOtpCacheClient(cacheServiceURL, nil)

	eventPublisher, err := client.NewOtpEventPublisher(natsURL, otpSubject)
	if err != nil {
		log.Printf("OTP NATS publisher disabled: %v", err)
		eventPublisher, _ = client.NewOtpEventPublisher("", otpSubject)
	}

	otpService := serviceimpl.NewOtpService(cacheClient, eventPublisher)
	mobileNorm := mobile.NewNormalizer(defaultCountryCode)
	otpHandler := handler.NewOtpHandler(otpService, mobileNorm)

	return &OtpContainer{
		OtpHandler: otpHandler,
	}
}
