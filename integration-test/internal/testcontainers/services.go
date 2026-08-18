package testcontainers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ServiceContainer struct {
	Name string
	testcontainers.Container
	Port nat.Port
}

func (s *ServiceContainer) BaseURL(ctx context.Context) (string, error) {
	host, err := s.Host(ctx)
	if err != nil {
		return "", err
	}

	mappedPort, err := s.MappedPort(ctx, s.Port)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%s", host, mappedPort.Port()), nil
}

func startService(ctx context.Context, net *testcontainers.DockerNetwork, cfg serviceConfig) (*ServiceContainer, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    RepoRoot(),
			Dockerfile: cfg.dockerfile,
			// Stable tags + KeepImage so parallel OTP packages reuse builds instead of
			// each package rebuilding with a random UUID tag (that caused 15m timeouts).
			Repo:      "ads-platform-it/" + cfg.name,
			Tag:       "local",
			KeepImage: true,
		},
		ExposedPorts: []string{string(cfg.port)},
		Env:          cfg.env,
		Networks:     []string{net.Name},
		NetworkAliases: map[string][]string{
			net.Name: {cfg.alias},
		},
		WaitingFor: wait.ForHTTP(cfg.healthPath).
			WithPort(cfg.port).
			WithStartupTimeout(5 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("start %s: %w", cfg.name, err)
	}

	return &ServiceContainer{
		Name:      cfg.name,
		Container: container,
		Port:      cfg.port,
	}, nil
}

type serviceConfig struct {
	name       string
	alias      string
	dockerfile string
	port       nat.Port
	healthPath string
	env        map[string]string
}

func StartNatsBroker(ctx context.Context, net *testcontainers.DockerNetwork) (*ServiceContainer, error) {
	return startService(ctx, net, serviceConfig{
		name:       "nats-message-broker",
		alias:      "nats-message-broker",
		dockerfile: "integration-test/docker/nats-message-broker.Dockerfile",
		port:       "8095/tcp",
		healthPath: "/health",
		env: map[string]string{
			"PORT":           "8095",
			"NATS_HOST":      "0.0.0.0",
			"NATS_PORT":      "4222",
			"NATS_HTTP_PORT": "8222",
		},
	})
}

func StartCacheService(ctx context.Context, net *testcontainers.DockerNetwork) (*ServiceContainer, error) {
	return startService(ctx, net, serviceConfig{
		name:       "ads-platform-cache-service",
		alias:      "ads-platform-cache-service",
		dockerfile: "integration-test/docker/ads-platform-cache-service.Dockerfile",
		port:       "8093/tcp",
		healthPath: "/health",
		env: map[string]string{
			"PORT":         "8093",
			"CDN_BASE_URL": "http://ads-platform-cdn:4000",
		},
	})
}

func StartCDN(ctx context.Context, net *testcontainers.DockerNetwork) (*ServiceContainer, error) {
	return startService(ctx, net, serviceConfig{
		name:       "ads-platform-cdn",
		alias:      "ads-platform-cdn",
		dockerfile: "integration-test/docker/ads-platform-cdn.Dockerfile",
		port:       "4000/tcp",
		healthPath: "/health",
		env: map[string]string{
			"PORT": "4000",
		},
	})
}

func StartNotification(ctx context.Context, net *testcontainers.DockerNetwork) (*ServiceContainer, error) {
	return startService(ctx, net, serviceConfig{
		name:       "ads-platform-notification",
		alias:      "ads-platform-notification",
		dockerfile: "integration-test/docker/ads-platform-notification.Dockerfile",
		port:       "8096/tcp",
		healthPath: "/health",
		env: map[string]string{
			"PORT":            "8096",
			"NATS_URL":        "nats://nats-message-broker:4222",
			"NATS_BROKER_URL": "http://nats-message-broker:8095",
			"OTP_SUBJECT":     "notifications.otp.send",
		},
	})
}

func StartBack(ctx context.Context, net *testcontainers.DockerNetwork, pg *PostgresContainer) (*ServiceContainer, error) {
	return startService(ctx, net, serviceConfig{
		name:       "ads-platform-back",
		alias:      "ads-platform-back",
		dockerfile: "integration-test/docker/ads-platform-back.Dockerfile",
		port:       "8092/tcp",
		healthPath: "/health",
		env: map[string]string{
			"APPLICATION_SERVER_PORT": "8092",
			"DATABASE_URL":            pg.DSN,
			"DATABASE_TYPE":           "postgres",
			"CACHE_SERVICE_URL":       "http://ads-platform-cache-service:8093",
			"NATS_URL":                "nats://nats-message-broker:4222",
			"NATS_BROKER_URL":         "http://nats-message-broker:8095",
			"OTP_SUBJECT":             "notifications.otp.send",
		},
	})
}

func StartBFF(ctx context.Context, net *testcontainers.DockerNetwork) (*ServiceContainer, error) {
	return startService(ctx, net, serviceConfig{
		name:       "ads-bff",
		alias:      "ads-bff",
		dockerfile: "integration-test/docker/ads-bff.Dockerfile",
		port:       "8097/tcp",
		healthPath: "/health",
		env: map[string]string{
			"PORT":                 "8097",
			"BACKEND_API_BASE_URL": "http://ads-platform-back:8092",
			"CACHE_SERVICE_URL":    "http://ads-platform-cache-service:8093",
		},
	})
}

func StartUI(ctx context.Context, net *testcontainers.DockerNetwork) (*ServiceContainer, error) {
	return startService(ctx, net, serviceConfig{
		name:       "ads-platform-ui",
		alias:      "ads-platform-ui",
		dockerfile: "integration-test/docker/ads-platform-ui.Dockerfile",
		port:       "8094/tcp",
		healthPath: "/health",
		env: map[string]string{
			"PORT":         "8094",
			"CDN_BASE_URL": "http://ads-platform-cdn:4000",
			"BFF_BASE_URL": "http://ads-bff:8097",
		},
	})
}

func WaitForHealthyURL(ctx context.Context, url string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	deadline := time.Now().Add(2 * time.Minute)

	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
			_ = body
		}

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("service not healthy at %s", url)
}

func Terminate(ctx context.Context, containers ...testcontainers.Container) {
	for _, container := range containers {
		if container == nil {
			continue
		}
		if err := container.Terminate(ctx); err != nil {
			log.Printf("terminate container: %v", err)
		}
	}
}
