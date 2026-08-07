//go:build integration

package e2e_happy_path

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestE2EHappyPath(t *testing.T) {
	_, backURL, cacheURL := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Send OTP
	sendStatus, sendResp, sendErr := otptest.SendOTP(ctx, backURL, otptest.TestMobile)
	if sendStatus != http.StatusOK || sendResp.Message != "otp_sent" {
		t.Fatalf("send failed: status=%d resp=%+v err=%+v", sendStatus, sendResp, sendErr)
	}

	// 2. OTP stored in cache
	otp, err := otptest.GetCachedOTP(ctx, cacheURL, otptest.TestMobile)
	if err != nil {
		t.Fatalf("cache lookup: %v", err)
	}

	// 3. Verify correct OTP
	verifyStatus, verifyResp, verifyErr := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, otp)
	if verifyStatus != http.StatusOK || !verifyResp.Verified {
		t.Fatalf("verify failed: status=%d resp=%+v err=%+v", verifyStatus, verifyResp, verifyErr)
	}

	// 4. Wrong OTP rejected
	wrongStatus, _, wrongErr := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, "000000")
	if wrongStatus != http.StatusUnauthorized || wrongErr.Error != "OTP_VERIFY_FAILED" {
		t.Fatalf("expected OTP_VERIFY_FAILED, status=%d err=%+v", wrongStatus, wrongErr)
	}
}
