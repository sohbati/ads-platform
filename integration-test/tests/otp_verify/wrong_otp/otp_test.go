//go:build integration

package wrong_otp

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestWrongOtp(t *testing.T) {
	_, backURL, cacheURL := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	otptest.MustSendAndGetOTP(t, ctx, backURL, cacheURL, otptest.TestMobile)

	status, _, errResp := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, "000000")
	if status != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", status)
	}
	if errResp.Error != "OTP_VERIFY_FAILED" {
		t.Fatalf("expected OTP_VERIFY_FAILED, got %q", errResp.Error)
	}
}
