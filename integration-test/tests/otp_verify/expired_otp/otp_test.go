//go:build integration

package expired_otp

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestExpiredOtp(t *testing.T) {
	t.Skip("OTP TTL is 500s; run manually or reduce cache TTL for automated expiry test")

	_, backURL, _ := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, _, errResp := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, "123456")
	if status != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", status)
	}
	if errResp.Error != "OTP_EXPIRED" {
		t.Fatalf("expected OTP_EXPIRED, got %q", errResp.Error)
	}
}
