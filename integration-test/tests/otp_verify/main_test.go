//go:build integration

package otp_verify

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	tc "integration-test/internal/testcontainers"
)

// Shared suite URLs — stack starts once for all OTP tests in this package.
var (
	backURL  string
	cacheURL string
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()

	stack, err := tc.StartOtpStack(ctx)
	if err != nil {
		panic("start otp stack: " + err.Error())
	}

	bu, err := stack.BackURL(ctx)
	if err != nil {
		stack.Terminate(context.Background())
		panic("back url: " + err.Error())
	}
	cu, err := stack.CacheURL(ctx)
	if err != nil {
		stack.Terminate(context.Background())
		panic("cache url: " + err.Error())
	}
	backURL, cacheURL = bu, cu

	code := m.Run()

	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cleanupCancel()
	stack.Terminate(cleanupCtx)
	os.Exit(code)
}
