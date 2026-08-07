//go:build integration

package valid_verify

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestValidVerify(t *testing.T) {
	_, backURL, cacheURL := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	otp := otptest.MustSendAndGetOTP(t, ctx, backURL, cacheURL, otptest.TestMobile)

	status, resp, errResp := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, otp)
	if status != http.StatusOK {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if !resp.Verified || resp.Message != "otp_verified" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
