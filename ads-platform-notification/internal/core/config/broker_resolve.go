package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type brokerHealth struct {
	BrokerURL  string `json:"brokerURL"`
	BrokerPort int    `json:"brokerPort"`
}

func resolveBrokerURL(directURL string, brokerHTTPURL string) (string, error) {
	if strings.TrimSpace(directURL) != "" {
		return strings.TrimSpace(directURL), nil
	}

	brokerHTTPURL = strings.TrimRight(strings.TrimSpace(brokerHTTPURL), "/")
	if brokerHTTPURL == "" {
		return "", fmt.Errorf("BROKER_URL is empty and BROKER_HTTP_URL is not set")
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(brokerHTTPURL + "/health")
	if err != nil {
		return "", fmt.Errorf("fetch broker health from %s: %w", brokerHTTPURL, err)
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

	return rewriteBrokerURL(health, hostFromHTTP(brokerHTTPURL))
}

func rewriteBrokerURL(health brokerHealth, host string) (string, error) {
	if strings.TrimSpace(health.BrokerURL) != "" {
		u, err := url.Parse(health.BrokerURL)
		if err != nil {
			return "", fmt.Errorf("parse brokerURL from health: %w", err)
		}
		if u.Scheme == "" {
			return "", fmt.Errorf("broker health brokerURL is missing a scheme")
		}
		port := u.Port()
		if health.BrokerPort > 0 {
			port = strconv.Itoa(health.BrokerPort)
		}
		if port == "" {
			return "", fmt.Errorf("broker health did not return a valid brokerPort")
		}
		u.Host = net.JoinHostPort(host, port)
		return u.String(), nil
	}

	if health.BrokerPort <= 0 {
		return "", fmt.Errorf("broker health did not return a valid brokerPort")
	}

	return fmt.Sprintf("%s://%s:%d", "nats", host, health.BrokerPort), nil
}

func hostFromHTTP(brokerHTTPURL string) string {
	if u, err := url.Parse(brokerHTTPURL); err == nil && u.Hostname() != "" {
		return u.Hostname()
	}
	return "127.0.0.1"
}
