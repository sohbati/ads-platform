//go:build integration

package resend_overwrites

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestResendOverwrites(t *testing.T) {
	_, backURL, cacheURL := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	firstOTP := otptest.MustSendAndGetOTP(t, ctx, backURL, cacheURL, otptest.TestMobile)
	secondOTP := otptest.MustSendAndGetOTP(t, ctx, backURL, cacheURL, otptest.TestMobile)

	if firstOTP == secondOTP {
		t.Log("warning: same OTP generated twice (possible but unlikely)")
	}

	cached, err := otptest.GetCachedOTP(ctx, cacheURL, otptest.TestMobile)
	if err != nil {
		t.Fatalf("get cached otp: %v", err)
	}
	if cached != secondOTP {
		t.Fatalf("cache should hold latest otp: cached=%q second=%q", cached, secondOTP)
	}

	status, _, errResp := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, firstOTP)
	if firstOTP != secondOTP && status == http.StatusOK {
		t.Fatalf("old otp should not verify after resend")
	}
	if firstOTP != secondOTP && errResp.Error != "OTP_VERIFY_FAILED" {
		t.Fatalf("expected OTP_VERIFY_FAILED for old otp, got %q", errResp.Error)
	}
}
