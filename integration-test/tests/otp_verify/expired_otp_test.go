//go:build integration

package otp_verify

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestExpiredOtp(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	otp := otptest.MustSendAndGetOTP(t, ctx, backURL, cacheURL, otptest.TestMobile)

	status, _, errResp := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, otp)
	if status != http.StatusOK {
		t.Fatalf("otp should verify during the wait window: status=%d error=%+v", status, errResp)
	}

	time.Sleep(2 * time.Second)

	status, _, errResp = otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, otp)
	if status != http.StatusNotFound {
		t.Fatalf("expected 404 after timer, got %d error=%+v", status, errResp)
	}
	if errResp.Error != "OTP_EXPIRED" {
		t.Fatalf("expected OTP_EXPIRED, got %q", errResp.Error)
	}
}
