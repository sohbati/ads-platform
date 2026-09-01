package impl

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"sync"
	"time"

	"ads-platform/internal/business/otp/client"
	"ads-platform/internal/business/otp/errorcode"
	"ads-platform/internal/business/otp/model"
	"ads-platform/internal/business/otp/service"
	"ads-platform/internal/core/exception"
)

const defaultResendAfter = 60 * time.Second

type otpService struct {
	cacheClient    client.OtpCacheClient
	eventPublisher client.OtpEventPublisher
	resendAfter    time.Duration
	now            func() time.Time
	locks          sync.Map
}

func NewOtpService(cacheClient client.OtpCacheClient, eventPublisher client.OtpEventPublisher, resendAfter time.Duration) service.OtpService {
	if resendAfter < 0 {
		resendAfter = defaultResendAfter
	}
	return &otpService{
		cacheClient:    cacheClient,
		eventPublisher: eventPublisher,
		resendAfter:    resendAfter,
		now:            time.Now,
	}
}

func (s *otpService) SendOTP(ctx context.Context, mobile string) (*model.SendOtpResponse, error) {
	if mobile == "" {
		return nil, exception.NewAppError(
			errorcode.ErrMobileEmpty.Code, errorcode.ErrMobileEmpty.HttpStatus, mobile)
	}

	unlock := s.lockMobile(mobile)
	defer unlock()

	if remaining := s.cooldownRemaining(ctx, mobile); remaining > 0 {
		return nil, exception.NewAppError(
			errorcode.ErrOTPResendWait.Code,
			errorcode.ErrOTPResendWait.HttpStatus,
			strconv.Itoa(remaining),
		)
	}

	otp, err := generateOTP()
	log.Printf("otp: %s", otp)
	if err != nil {
		return nil, fmt.Errorf("generate otp: %w", err)
	}

	if err := s.cacheClient.StoreOTP(ctx, otpCacheKey(mobile), otp); err != nil {
		return nil, exception.NewAppError(
			errorcode.ErrCacheUnavailable.Code, errorcode.ErrCacheUnavailable.HttpStatus, mobile).WithCause(err)
	}

	if s.resendAfter > 0 {
		until := s.now().Add(s.resendAfter).Unix()
		if err := s.cacheClient.StoreOTP(ctx, cooldownCacheKey(mobile), strconv.FormatInt(until, 10)); err != nil {
			return nil, exception.NewAppError(
				errorcode.ErrCacheUnavailable.Code, errorcode.ErrCacheUnavailable.HttpStatus, mobile).WithCause(err)
		}
	}

	if err := s.eventPublisher.PublishOtpEvent(ctx, mobile, otp); err != nil {
		log.Printf("failed to publish otp event for mobile=%s: %v", mobile, err)
	}

	return &model.SendOtpResponse{
		Message:            "otp_sent",
		ResendAfterSeconds: int(s.resendAfter / time.Second),
	}, nil
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

	if s.otpExpired(ctx, mobile) {
		return nil, exception.NewAppError(
			errorcode.ErrOTPExpired.Code, errorcode.ErrOTPExpired.HttpStatus, mobile)
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

func (s *otpService) cooldownRemaining(ctx context.Context, mobile string) int {
	if s.resendAfter <= 0 {
		return 0
	}
	raw, err := s.cacheClient.GetOTP(ctx, cooldownCacheKey(mobile))
	if err != nil {
		return 0
	}
	untilUnix, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0
	}
	remaining := time.Unix(untilUnix, 0).Sub(s.now())
	if remaining <= 0 {
		return 0
	}
	secs := int(math.Ceil(remaining.Seconds()))
	if secs < 1 {
		return 1
	}
	return secs
}

func (s *otpService) otpExpired(ctx context.Context, mobile string) bool {
	if s.resendAfter <= 0 {
		return false
	}
	raw, err := s.cacheClient.GetOTP(ctx, cooldownCacheKey(mobile))
	if err != nil {
		return true
	}
	untilUnix, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return true
	}
	return !s.now().Before(time.Unix(untilUnix, 0))
}

func (s *otpService) lockMobile(mobile string) func() {
	v, _ := s.locks.LoadOrStore(mobile, &sync.Mutex{})
	mu := v.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

func otpCacheKey(mobile string) string {
	return "otp:" + mobile
}

func cooldownCacheKey(mobile string) string {
	return "otp-cooldown:" + mobile
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
