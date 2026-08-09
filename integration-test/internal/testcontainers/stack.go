package testcontainers

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

type Stack struct {
	Network      *testcontainers.DockerNetwork
	Postgres     *PostgresContainer
	NatsBroker   *ServiceContainer
	Cache        *ServiceContainer
	CDN          *ServiceContainer
	Notification *ServiceContainer
	Back         *ServiceContainer
	BFF          *ServiceContainer
	UI           *ServiceContainer
}

func StartStack(ctx context.Context, t *testing.T) (*Stack, error) {
	t.Helper()

	net, err := CreateNetwork(ctx)
	if err != nil {
		return nil, fmt.Errorf("create network: %w", err)
	}

	pg, err := StartPostgres(ctx, t, net)
	if err != nil {
		_ = net.Remove(ctx)
		return nil, err
	}

	stack := &Stack{
		Network:  net,
		Postgres: pg,
	}

	start := func(name string, fn func() (*ServiceContainer, error)) (*ServiceContainer, error) {
		svc, err := fn()
		if err != nil {
			stack.Terminate(ctx, t)
			return nil, fmt.Errorf("start %s: %w", name, err)
		}
		return svc, nil
	}

	stack.NatsBroker, err = start("nats-message-broker", func() (*ServiceContainer, error) {
		return StartNatsBroker(ctx, t, net)
	})
	if err != nil {
		return nil, err
	}

	stack.Cache, err = start("ads-platform-cache-service", func() (*ServiceContainer, error) {
		return StartCacheService(ctx, t, net)
	})
	if err != nil {
		return nil, err
	}

	stack.CDN, err = start("ads-platform-cdn", func() (*ServiceContainer, error) {
		return StartCDN(ctx, t, net)
	})
	if err != nil {
		return nil, err
	}

	stack.Notification, err = start("ads-platform-notification", func() (*ServiceContainer, error) {
		return StartNotification(ctx, t, net)
	})
	if err != nil {
		return nil, err
	}

	stack.Back, err = start("ads-platform-back", func() (*ServiceContainer, error) {
		return StartBack(ctx, t, net, pg)
	})
	if err != nil {
		return nil, err
	}

	stack.BFF, err = start("ads-bff", func() (*ServiceContainer, error) {
		return StartBFF(ctx, t, net)
	})
	if err != nil {
		return nil, err
	}

	stack.UI, err = start("ads-platform-ui", func() (*ServiceContainer, error) {
		return StartUI(ctx, t, net)
	})
	if err != nil {
		return nil, err
	}

	return stack, nil
}

func (s *Stack) Terminate(ctx context.Context, t *testing.T) {
	t.Helper()

	Terminate(ctx, t,
		serviceContainer(s.UI),
		serviceContainer(s.BFF),
		serviceContainer(s.Back),
		serviceContainer(s.Notification),
		serviceContainer(s.CDN),
		serviceContainer(s.Cache),
		serviceContainer(s.NatsBroker),
	)

	if s.Postgres != nil && s.Postgres.Container != nil {
		if err := s.Postgres.Container.Terminate(ctx); err != nil {
			t.Logf("terminate postgres: %v", err)
		}
	}

	if s.Network != nil {
		if err := s.Network.Remove(ctx); err != nil {
			t.Logf("remove network: %v", err)
		}
	}
}

func serviceContainer(svc *ServiceContainer) testcontainers.Container {
	if svc == nil {
		return nil
	}
	return svc.Container
}

func (s *Stack) URLs(ctx context.Context) (map[string]string, error) {
	urls := make(map[string]string)
	services := map[string]*ServiceContainer{
		"ads-platform-back":          s.Back,
		"ads-bff":                    s.BFF,
		"ads-platform-cache-service": s.Cache,
		"ads-platform-cdn":           s.CDN,
		"ads-platform-ui":            s.UI,
		"nats-message-broker":        s.NatsBroker,
		"ads-platform-notification":  s.Notification,
	}

	for name, svc := range services {
		url, err := svc.BaseURL(ctx)
		if err != nil {
			return nil, fmt.Errorf("%s base url: %w", name, err)
		}
		urls[name] = url
	}

	return urls, nil
}
