package container

import (
	"ads-platform/internal/business/otp/client"
	"ads-platform/internal/business/otp/handler"
	serviceimpl "ads-platform/internal/business/otp/service/impl"
)

type OtpContainer struct {
	OtpHandler *handler.OtpHandler
}

func NewOtpContainer(cacheServiceURL string) *OtpContainer {
	cacheClient := client.NewOtpCacheClient(cacheServiceURL, nil)
	otpService := serviceimpl.NewOtpService(cacheClient)
	otpHandler := handler.NewOtpHandler(otpService)

	return &OtpContainer{
		OtpHandler: otpHandler,
	}
}
