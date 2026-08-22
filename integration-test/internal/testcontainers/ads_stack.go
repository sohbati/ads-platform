package testcontainers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// AdsAPIStack is postgres + CDN + cache + MinIO + ads-platform-back for POST /ads.
type AdsAPIStack struct {
	Network  *testcontainers.DockerNetwork
	Postgres *PostgresContainer
	Minio    *MinioContainer
	CDN      *ServiceContainer
	Cache    *ServiceContainer
	Back     *ServiceContainer
}

func StartAdsAPIStack(ctx context.Context) (*AdsAPIStack, error) {
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

	stack := &AdsAPIStack{Network: net, Postgres: pg}

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

	stack.Minio, err = StartMinio(ctx, net)
	if err != nil {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		stack.Terminate(cleanupCtx)
		return nil, fmt.Errorf("start minio: %w", err)
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

	stack.Back, err = start("ads-platform-back", func() (*ServiceContainer, error) {
		return StartBackWithEnv(ctx, net, pg, map[string]string{
			"MINIO_ENDPOINT":        stack.Minio.Endpoint,
			"MINIO_ACCESS_KEY":      stack.Minio.AccessKey,
			"MINIO_SECRET_KEY":      stack.Minio.SecretKey,
			"MINIO_BUCKET":          stack.Minio.Bucket,
			"MINIO_USE_SSL":         "false",
			"MINIO_PUBLIC_URL":      "http://" + stack.Minio.Endpoint,
			"ADS_MAX_PICTURES":      "2",
			"ADS_MAX_PICTURE_BYTES": "10485760",
		})
	})
	if err != nil {
		return nil, err
	}

	return stack, nil
}

func (s *AdsAPIStack) Terminate(ctx context.Context) {
	Terminate(ctx,
		serviceContainer(s.Back),
		serviceContainer(s.Cache),
		serviceContainer(s.CDN),
	)
	if s.Minio != nil && s.Minio.Container != nil {
		if err := s.Minio.Container.Terminate(ctx); err != nil {
			log.Printf("terminate minio: %v", err)
		}
	}
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

func (s *AdsAPIStack) BackURL(ctx context.Context) (string, error) {
	return s.Back.BaseURL(ctx)
}
