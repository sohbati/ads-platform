//go:build integration

package missing_body

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestMissingBody(t *testing.T) {
	_, backURL, _ := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, _, errResp := otptest.VerifyOTPRaw(ctx, backURL, otptest.TestMobile, nil)
	if status != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", status)
	}
	if errResp.Error != "BAD_REQUEST" {
		t.Fatalf("expected BAD_REQUEST, got %q", errResp.Error)
	}
}
