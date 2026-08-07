//go:build integration

package empty_mobile_send

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestEmptyMobileSend(t *testing.T) {
	_, backURL, _ := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	status, _, errResp := otptest.SendOTP(ctx, backURL, "")
	if status != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", status)
	}
	if errResp.Error != "MOBILE_EMPTY" {
		t.Fatalf("expected MOBILE_EMPTY, got %q", errResp.Error)
	}
}
