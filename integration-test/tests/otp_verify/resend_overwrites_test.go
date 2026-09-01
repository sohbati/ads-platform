//go:build integration

package otp_verify

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestResendBlockedUntilWaitElapses(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	firstOTP := otptest.MustSendAndGetOTP(t, ctx, backURL, cacheURL, otptest.TestMobile)

	status, _, errResp := otptest.SendOTP(ctx, backURL, otptest.TestMobile)
	if status != http.StatusTooManyRequests {
		t.Fatalf("immediate resend status=%d error=%+v, want 429", status, errResp)
	}
	if errResp.Error != "OTP_RESEND_WAIT" {
		t.Fatalf("expected OTP_RESEND_WAIT, got %q", errResp.Error)
	}
	cached, err := otptest.GetCachedOTP(ctx, cacheURL, otptest.TestMobile)
	if err != nil {
		t.Fatalf("get cached otp: %v", err)
	}
	if cached != firstOTP {
		t.Fatalf("otp must not change during wait: cached=%q first=%q", cached, firstOTP)
	}

	secondOTP := otptest.MustSendAndGetOTP(t, ctx, backURL, cacheURL, otptest.TestMobile)
	if firstOTP == secondOTP {
		t.Log("warning: same OTP generated twice (possible but unlikely)")
	}

	status, _, errResp = otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, firstOTP)
	if firstOTP != secondOTP && status == http.StatusOK {
		t.Fatalf("old otp should not verify after resend")
	}
	if firstOTP != secondOTP && errResp.Error != "OTP_VERIFY_FAILED" {
		t.Fatalf("expected OTP_VERIFY_FAILED for old otp, got %q", errResp.Error)
	}
}
