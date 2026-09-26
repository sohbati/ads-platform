package testcontainers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// IAMStack is postgres + ads-platform-back for IAM HTTP APIs.
type IAMStack struct {
	Network  *testcontainers.DockerNetwork
	Postgres *PostgresContainer
	Back     *ServiceContainer
}

func StartIAMStack(ctx context.Context) (*IAMStack, error) {
	net, err := CreateNetwork(ctx)
	if err != nil {
		return nil, fmt.Errorf("create network: %w", err)
	}

	pg, err := StartPostgres(ctx, net)
	if err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		_ = net.Remove(cleanupCtx)
		return nil, err
	}

	stack := &IAMStack{Network: net, Postgres: pg}

	stack.Back, err = StartBack(ctx, net, pg)
	if err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		stack.Terminate(cleanupCtx)
		return nil, fmt.Errorf("start ads-platform-back: %w", err)
	}

	return stack, nil
}

func (s *IAMStack) Terminate(ctx context.Context) {
	Terminate(ctx, serviceContainer(s.Back))
	if s.Postgres != nil && s.Postgres.Container != nil {
		if err := s.Postgres.Container.Terminate(ctx); err != nil {
			log.Printf("terminate postgres: %v", err)
		}
	}
	if s.Network != nil {
		if err := s.Network.Remove(ctx); err != nil {
			log.Printf("remove network: %v", err)
		}
	}
}

func (s *IAMStack) BackURL(ctx context.Context) (string, error) {
	return s.Back.BaseURL(ctx)
}
