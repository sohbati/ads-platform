package bff

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client calls ads-bff, which proxies to ads-platform-back.
type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		http:    httpClient,
	}
}

// Do forwards an HTTP request to ads-bff and returns the upstream status and body.
func (c *Client) Do(ctx context.Context, method, path string, body []byte, contentType string) (int, []byte, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return 0, nil, fmt.Errorf("bff: build request %s %s: %w", method, path, err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("bff: %s %s: %w", method, url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return 0, nil, fmt.Errorf("bff: read %s: %w", url, err)
	}

	return resp.StatusCode, data, nil
}

func (c *Client) Get(ctx context.Context, path string) (int, []byte, error) {
	return c.Do(ctx, http.MethodGet, path, nil, "")
}

func (c *Client) PostJSON(ctx context.Context, path string, body []byte) (int, []byte, error) {
	return c.Do(ctx, http.MethodPost, path, body, "application/json")
}
