package testcontainers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// CityCacheStack is CDN + cache-service for city/category cache APIs.
type CityCacheStack struct {
	Network *testcontainers.DockerNetwork
	CDN     *ServiceContainer
	Cache   *ServiceContainer
}

func StartCityCacheStack(ctx context.Context) (*CityCacheStack, error) {
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
			stack.Terminate(cleanupCtx)
			return nil, fmt.Errorf("start %s: %w", name, err)
		}
		return svc, nil
	}

	stack.CDN, err = start("ads-platform-cdn", func() (*ServiceContainer, error) {
		return StartCDN(ctx, net)
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

	return stack, nil
}

func (s *CityCacheStack) Terminate(ctx context.Context) {
	Terminate(ctx,
		serviceContainer(s.Cache),
		serviceContainer(s.CDN),
	)

	if s.Network != nil {
		if err := s.Network.Remove(ctx); err != nil {
			log.Printf("remove network: %v", err)
		}
	}
}

func (s *CityCacheStack) CacheURL(ctx context.Context) (string, error) {
	return s.Cache.BaseURL(ctx)
}

func (s *CityCacheStack) CDNURL(ctx context.Context) (string, error) {
	return s.CDN.BaseURL(ctx)
}
