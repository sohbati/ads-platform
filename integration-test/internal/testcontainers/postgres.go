package testcontainers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	networkName = "ads-platform-integration"

	pgUser     = "ads_platform_user"
	pgPassword = "ads_platform_user"
	pgDatabase = "ads_platform"
)

type PostgresContainer struct {
	Container *postgres.PostgresContainer
	// DSN is reachable from other containers on the shared Docker network.
	DSN string
	// HostDSN is reachable from the test process on the host.
	HostDSN string
}

func StartPostgres(ctx context.Context, t *testing.T, net *testcontainers.DockerNetwork) (*PostgresContainer, error) {
	t.Helper()

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(pgDatabase),
		postgres.WithUsername(pgUser),
		postgres.WithPassword(pgPassword),
		postgres.WithInitScripts(PostgresInitScript()),
		network.WithNetwork([]string{networkName}, net),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(2*time.Minute),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("start postgres: %w", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?search_path=ads_platform_schema&sslmode=disable",
		pgUser, pgPassword, networkName, pgDatabase,
	)

	hostDSN, err := container.ConnectionString(ctx, "sslmode=disable", "search_path=ads_platform_schema")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("postgres host connection string: %w", err)
	}

	return &PostgresContainer{
		Container: container,
		DSN:       dsn,
		HostDSN:   hostDSN,
	}, nil
}

func CreateNetwork(ctx context.Context) (*testcontainers.DockerNetwork, error) {
	return network.New(ctx)
}
