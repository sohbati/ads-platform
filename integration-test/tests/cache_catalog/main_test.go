//go:build integration

package cache_catalog

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	tc "integration-test/internal/testcontainers"
)

// Shared CDN + cache-service stack for all city/category cache API tests.
var cacheURL string

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()

	stack, err := tc.StartCityCacheStack(ctx)
	if err != nil {
		panic("start cache catalog stack: " + err.Error())
	}

	url, err := stack.CacheURL(ctx)
	if err != nil {
		stack.Terminate(context.Background())
		panic("cache url: " + err.Error())
	}
	cacheURL = url

	code := m.Run()

	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cleanupCancel()
	stack.Terminate(cleanupCtx)
	os.Exit(code)
}
