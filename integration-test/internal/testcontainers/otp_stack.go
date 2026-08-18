package testcontainers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// OtpStack contains services required for OTP send/verify integration tests.
type OtpStack struct {
	Network      *testcontainers.DockerNetwork
	Postgres     *PostgresContainer
	NatsBroker   *ServiceContainer
	Cache        *ServiceContainer
	Notification *ServiceContainer
	Back         *ServiceContainer
}

func StartOtpStack(ctx context.Context) (*OtpStack, error) {
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

	stack := &OtpStack{
		Network:  net,
		Postgres: pg,
	}

	start := func(name string, fn func() (*ServiceContainer, error)) (*ServiceContainer, error) {
		svc, err := fn()
		if err != nil {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			stack.Terminate(cleanupCtx)
			return nil, fmt.Errorf("start %s: %w", name, err)
		}
		return svc, nil
	}

	stack.NatsBroker, err = start("nats-message-broker", func() (*ServiceContainer, error) {
		return StartNatsBroker(ctx, net)
	})
	if err != nil {
		return nil, err
	}

	stack.Cache, err = start("ads-platform-cache-service", func() (*ServiceContainer, error) {
		return StartCacheService(ctx, net)
	})
	if err != nil {
		return nil, err
	}

	stack.Notification, err = start("ads-platform-notification", func() (*ServiceContainer, error) {
		return StartNotification(ctx, net)
	})
	if err != nil {
		return nil, err
	}

	stack.Back, err = start("ads-platform-back", func() (*ServiceContainer, error) {
		return StartBack(ctx, net, pg)
	})
	if err != nil {
		return nil, err
	}

	return stack, nil
}

func (s *OtpStack) Terminate(ctx context.Context) {
	Terminate(ctx,
		serviceContainer(s.Back),
		serviceContainer(s.Notification),
		serviceContainer(s.Cache),
		serviceContainer(s.NatsBroker),
	)

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

func (s *OtpStack) BackURL(ctx context.Context) (string, error) {
	return s.Back.BaseURL(ctx)
}

func (s *OtpStack) CacheURL(ctx context.Context) (string, error) {
	return s.Cache.BaseURL(ctx)
}
