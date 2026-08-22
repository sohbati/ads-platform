package testcontainers

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	minioUser     = "minioadmin"
	minioPassword = "minioadmin123"
	minioBucket   = "ads-media"
	minioAlias    = "minio"
	minioAPIPort  = "9000/tcp"
)

type MinioContainer struct {
	Container testcontainers.Container
	AccessKey string
	SecretKey string
	Bucket    string
	// Endpoint is reachable from other containers on the shared network (minio:9000).
	Endpoint string
	// HostEndpoint is reachable from the test process (host:mappedPort).
	HostEndpoint string
}

func StartMinio(ctx context.Context, net *testcontainers.DockerNetwork) (*MinioContainer, error) {
	port := nat.Port(minioAPIPort)
	req := testcontainers.ContainerRequest{
		Image:        "minio/minio:latest",
		ExposedPorts: []string{string(port)},
		Env: map[string]string{
			"MINIO_ROOT_USER":     minioUser,
			"MINIO_ROOT_PASSWORD": minioPassword,
		},
		Cmd:      []string{"server", "/data"},
		Networks: []string{net.Name},
		NetworkAliases: map[string][]string{
			net.Name: {minioAlias},
		},
		WaitingFor: wait.ForHTTP("/minio/health/live").
			WithPort(port).
			WithStartupTimeout(2 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start minio: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("minio host: %w", err)
	}
	mapped, err := container.MappedPort(ctx, port)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, fmt.Errorf("minio mapped port: %w", err)
	}

	return &MinioContainer{
		Container:    container,
		AccessKey:    minioUser,
		SecretKey:    minioPassword,
		Bucket:       minioBucket,
		Endpoint:     minioAlias + ":9000",
		HostEndpoint: fmt.Sprintf("%s:%s", host, mapped.Port()),
	}, nil
}
