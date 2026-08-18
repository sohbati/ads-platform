package testcontainers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// CityCacheStack is CDN + cache-service for city lookup APIs.
type CityCacheStack struct {
	Network *testcontainers.DockerNetwork
	CDN     *ServiceContainer
	Cache   *ServiceContainer
}

func StartCityCacheStack(ctx context.Context, t *testing.T) (*CityCacheStack, error) {
	t.Helper()

	net, err := CreateNetwork(ctx)
	if err != nil {
		return nil, fmt.Errorf("create network: %w", err)
	}

	stack := &CityCacheStack{Network: net}

	start := func(name string, fn func() (*ServiceContainer, error)) (*ServiceContainer, error) {
		svc, err := fn()
		if err != nil {
			cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			stack.Terminate(cleanupCtx, t)
			return nil, fmt.Errorf("start %s: %w", name, err)
		}
		return svc, nil
	}

	stack.CDN, err = start("ads-platform-cdn", func() (*ServiceContainer, error) {
		return StartCDN(ctx, t, net)
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

	return stack, nil
}

func (s *CityCacheStack) Terminate(ctx context.Context, t *testing.T) {
	t.Helper()

	Terminate(ctx, t,
		serviceContainer(s.Cache),
		serviceContainer(s.CDN),
	)

	if s.Network != nil {
		if err := s.Network.Remove(ctx); err != nil {
			t.Logf("remove network: %v", err)
		}
	}
}

func (s *CityCacheStack) CacheURL(ctx context.Context) (string, error) {
	return s.Cache.BaseURL(ctx)
}

func (s *CityCacheStack) CDNURL(ctx context.Context) (string, error) {
	return s.CDN.BaseURL(ctx)
}
