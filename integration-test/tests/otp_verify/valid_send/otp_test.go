//go:build integration

package valid_send

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestValidSend(t *testing.T) {
	_, backURL, cacheURL := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, resp, errResp := otptest.SendOTP(ctx, backURL, otptest.TestMobile)
	if status != http.StatusOK {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if resp.Message != "otp_sent" {
		t.Fatalf("expected otp_sent, got %q", resp.Message)
	}

	otp, err := otptest.GetCachedOTP(ctx, cacheURL, otptest.TestMobile)
	if err != nil {
		t.Fatalf("cached otp: %v", err)
	}
	if len(otp) != 6 {
		t.Fatalf("expected 6-digit otp, got %q", otp)
	}
}
