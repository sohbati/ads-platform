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

// Client calls ads-bff.
type Client struct {
	baseURL string
	http    *http.Client
}

// HTTPResult includes upstream headers (e.g. Set-Cookie).
type HTTPResult struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 60 * time.Second}
	}
	return &Client{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		http:    httpClient,
	}
}

func (c *Client) Do(ctx context.Context, method, path string, body []byte, contentType string, cookies string) (*HTTPResult, error) {
	url := c.baseURL + path

	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("bff: build request %s %s: %w", method, path, err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if cookies != "" {
		req.Header.Set("Cookie", cookies)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bff: %s %s: %w", method, url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return nil, fmt.Errorf("bff: read %s: %w", url, err)
	}

	return &HTTPResult{
		StatusCode: resp.StatusCode,
		Body:       data,
		Header:     resp.Header.Clone(),
	}, nil
}

func (c *Client) Get(ctx context.Context, path, cookies string) (*HTTPResult, error) {
	return c.Do(ctx, http.MethodGet, path, nil, "", cookies)
}

func (c *Client) PostJSON(ctx context.Context, path string, body []byte, cookies string) (*HTTPResult, error) {
	return c.Do(ctx, http.MethodPost, path, body, "application/json", cookies)
}

func ForwardResponse(c http.ResponseWriter, result *HTTPResult) {
	for _, cookie := range result.Header["Set-Cookie"] {
		c.Header().Add("Set-Cookie", cookie)
	}
	contentType := result.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}
	c.WriteHeader(result.StatusCode)
	_, _ = c.Write(result.Body)
}

func RequestCookies(r *http.Request) string {
	return r.Header.Get("Cookie")
}
