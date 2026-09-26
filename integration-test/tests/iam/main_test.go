//go:build integration

package iam

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	tc "integration-test/internal/testcontainers"
)

var (
	backURL   string
	pgHostDSN string
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()

	stack, err := tc.StartIAMStack(ctx)
	if err != nil {
		panic("start iam stack: " + err.Error())
	}

	url, err := stack.BackURL(ctx)
	if err != nil {
		stack.Terminate(context.Background())
		panic("back url: " + err.Error())
	}
	backURL = url
	pgHostDSN = stack.Postgres.HostDSN

	code := m.Run()

	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cleanupCancel()
	stack.Terminate(cleanupCtx)
	os.Exit(code)
}
