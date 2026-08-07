//go:build integration

package wrong_mobile

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestWrongMobile(t *testing.T) {
	_, backURL, cacheURL := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	otp := otptest.MustSendAndGetOTP(t, ctx, backURL, cacheURL, otptest.TestMobile)

	status, _, errResp := otptest.VerifyOTP(ctx, backURL, otptest.OtherMobile, otp)
	if status != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", status)
	}
	if errResp.Error != "OTP_EXPIRED" {
		t.Fatalf("expected OTP_EXPIRED, got %q", errResp.Error)
	}
}
