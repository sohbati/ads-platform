//go:build integration

package verify_twice

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/otptest"
)

func TestVerifyTwice(t *testing.T) {
	_, backURL, cacheURL := otptest.SetupOtpStack(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	otp := otptest.MustSendAndGetOTP(t, ctx, backURL, cacheURL, otptest.TestMobile)

	status1, resp1, err1 := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, otp)
	if status1 != http.StatusOK || !resp1.Verified {
		t.Fatalf("first verify failed: status=%d resp=%+v err=%+v", status1, resp1, err1)
	}

	status2, resp2, err2 := otptest.VerifyOTP(ctx, backURL, otptest.TestMobile, otp)
	if status2 != http.StatusOK || !resp2.Verified {
		t.Fatalf("second verify failed: status=%d resp=%+v err=%+v", status2, resp2, err2)
	}
}
