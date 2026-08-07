package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type brokerHealth struct {
	NatsPort int `json:"natsPort"`
}

func resolveNatsURL(natsURL string, brokerURL string) (string, error) {
	if strings.TrimSpace(natsURL) != "" {
		return strings.TrimSpace(natsURL), nil
	}

	brokerURL = strings.TrimRight(strings.TrimSpace(brokerURL), "/")
	if brokerURL == "" {
		return "", fmt.Errorf("NATS_URL is empty and NATS_BROKER_URL is not set")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(brokerURL + "/health")
	if err != nil {
		return "", fmt.Errorf("fetch broker health from %s: %w", brokerURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("broker health returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<10))
	if err != nil {
		return "", fmt.Errorf("read broker health: %w", err)
	}

	var health brokerHealth
	if err := json.Unmarshal(body, &health); err != nil {
		return "", fmt.Errorf("parse broker health: %w", err)
	}

	if health.NatsPort <= 0 {
		return "", fmt.Errorf("broker health did not return a valid natsPort")
	}

	host := "127.0.0.1"
	if strings.HasPrefix(brokerURL, "http://") {
		rest := strings.TrimPrefix(brokerURL, "http://")
		if idx := strings.Index(rest, ":"); idx > 0 {
			host = rest[:idx]
		} else if rest != "" {
			host = rest
		}
	}

	return fmt.Sprintf("nats://%s:%d", host, health.NatsPort), nil
}
