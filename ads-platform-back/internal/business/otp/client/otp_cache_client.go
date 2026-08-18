package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"ads-platform/internal/business/otp/model"
)

type OtpCacheClient interface {
	StoreOTP(ctx context.Context, key string, otp string) error
	GetOTP(ctx context.Context, key string) (string, error)
}

type otpCacheClient struct {
	baseURL string
	http    *http.Client
}

func NewOtpCacheClient(baseURL string, httpClient *http.Client) OtpCacheClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &otpCacheClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		http:    httpClient,
	}
}

func (c *otpCacheClient) StoreOTP(ctx context.Context, key string, otp string) error {
	body, err := json.Marshal(model.CacheOtpRequest{Otp: otp})
	if err != nil {
		return fmt.Errorf("cache-service: marshal otp request: %w", err)
	}

	storeURL := fmt.Sprintf("%s/api/v1/caches/otp/%s", c.baseURL, url.PathEscape(key))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, storeURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cache-service: build store request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("cache-service: POST %s: %w", storeURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cache-service: POST %s: status %d", storeURL, resp.StatusCode)
	}

	return nil
}

func (c *otpCacheClient) GetOTP(ctx context.Context, key string) (string, error) {
	getURL := fmt.Sprintf("%s/api/v1/caches/otp/%s", c.baseURL, url.PathEscape(key))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return "", fmt.Errorf("cache-service: build get request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("cache-service: GET %s: %w", getURL, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<10))
	if err != nil {
		return "", fmt.Errorf("cache-service: read %s: %w", getURL, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", &cacheNotFoundError{statusCode: resp.StatusCode, body: string(data)}
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cache-service: GET %s: status %d", getURL, resp.StatusCode)
	}

	var otp string
	if err := json.Unmarshal(data, &otp); err != nil {
		return "", fmt.Errorf("cache-service: parse %s: %w", getURL, err)
	}

	return otp, nil
}

type cacheNotFoundError struct {
	statusCode int
	body       string
}

func (e *cacheNotFoundError) Error() string {
	return fmt.Sprintf("cache-service: otp not found (status %d)", e.statusCode)
}

func IsCacheNotFound(err error) bool {
	_, ok := err.(*cacheNotFoundError)
	return ok
}
