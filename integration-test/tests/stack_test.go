//go:build integration

package integration_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	tc "integration-test/internal/testcontainers"
)

func TestStackHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()

	stack, err := tc.StartStack(ctx, t)
	if err != nil {
		t.Fatalf("start stack: %v", err)
	}
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cleanupCancel()
		stack.Terminate(cleanupCtx, t)
	})

	urls, err := stack.URLs(ctx)
	if err != nil {
		t.Fatalf("resolve service urls: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	for name, url := range urls {
		t.Run(name, func(t *testing.T) {
			resp, err := client.Get(url + "/health")
			if err != nil {
				t.Fatalf("GET %s/health: %v", url, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("GET %s/health: status %d", url, resp.StatusCode)
			}
		})
	}
}

func TestPostgresConnection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	net, err := tc.CreateNetwork(ctx)
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	t.Cleanup(func() {
		_ = net.Remove(context.Background())
	})

	pg, err := tc.StartPostgres(ctx, t, net)
	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}
	t.Cleanup(func() {
		_ = pg.Container.Terminate(context.Background())
	})

	if pg.DSN == "" {
		t.Fatal("expected postgres dsn")
	}
}
