package impl

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"ads-platform/internal/business/otp/client"
	"ads-platform/internal/business/otp/errorcode"
	"ads-platform/internal/business/otp/model"
	"ads-platform/internal/business/otp/service"
	"ads-platform/internal/core/exception"
)

type otpService struct {
	cacheClient client.OtpCacheClient
}

func NewOtpService(cacheClient client.OtpCacheClient) service.OtpService {
	return &otpService{cacheClient: cacheClient}
}

func (s *otpService) SendOTP(ctx context.Context, mobile string) (*model.SendOtpResponse, error) {
	if mobile == "" {
		return nil, exception.NewAppError(
			errorcode.ErrMobileEmpty.Code, errorcode.ErrMobileEmpty.HttpStatus, mobile)
	}

	otp, err := generateOTP()
	if err != nil {
		return nil, fmt.Errorf("generate otp: %w", err)
	}

	if err := s.cacheClient.StoreOTP(ctx, otpCacheKey(mobile), otp); err != nil {
		return nil, exception.NewAppError(
			errorcode.ErrCacheUnavailable.Code, errorcode.ErrCacheUnavailable.HttpStatus, mobile).WithCause(err)
	}

	return &model.SendOtpResponse{Message: "otp_sent"}, nil
}

func (s *otpService) VerifyOTP(ctx context.Context, mobile string, otp string) (*model.VerifyOtpResponse, error) {
	if mobile == "" {
		return nil, exception.NewAppError(
			errorcode.ErrMobileEmpty.Code, errorcode.ErrMobileEmpty.HttpStatus, mobile)
	}

	if otp == "" {
		return nil, exception.NewAppError(
			errorcode.ErrInvalidOTP.Code, errorcode.ErrInvalidOTP.HttpStatus, mobile)
	}

	storedOTP, err := s.cacheClient.GetOTP(ctx, otpCacheKey(mobile))
	if err != nil {
		if client.IsCacheNotFound(err) {
			return nil, exception.NewAppError(
				errorcode.ErrOTPExpired.Code, errorcode.ErrOTPExpired.HttpStatus, mobile).WithCause(err)
		}
		return nil, exception.NewAppError(
			errorcode.ErrCacheUnavailable.Code, errorcode.ErrCacheUnavailable.HttpStatus, mobile).WithCause(err)
	}

	if storedOTP != otp {
		return nil, exception.NewAppError(
			errorcode.ErrOTPVerifyFailed.Code, errorcode.ErrOTPVerifyFailed.HttpStatus, mobile)
	}

	return &model.VerifyOtpResponse{
		Verified: true,
		Message:  "otp_verified",
	}, nil
}

func otpCacheKey(mobile string) string {
	return "otp:" + mobile
}

func generateOTP() (string, error) {
	const digits = "0123456789"
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[n.Int64()]
	}
	return string(code), nil
}
