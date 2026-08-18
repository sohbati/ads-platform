//go:build integration

package otp_verify

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestNonNumericOtp(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, _, errResp := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, "abc123")
	if status != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", status)
	}
	if errResp.Error != "BAD_REQUEST" {
		t.Fatalf("expected BAD_REQUEST, got %q", errResp.Error)
	}
}
