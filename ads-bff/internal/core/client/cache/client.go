package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &Client{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		http:    httpClient,
	}
}

type storeSessionRequest struct {
	Data string `json:"data"`
}

func (c *Client) StoreSession(ctx context.Context, key, data string) error {
	body, err := json.Marshal(storeSessionRequest{Data: data})
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/caches/session/%s", c.baseURL, key)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<10))
		return fmt.Errorf("store session status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (c *Client) GetSession(ctx context.Context, key string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/caches/session/%s", c.baseURL, key)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrSessionNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get session status %d", resp.StatusCode)
	}

	var sessionData string
	if err := json.Unmarshal(data, &sessionData); err != nil {
		return "", fmt.Errorf("parse session data: %w", err)
	}
	return sessionData, nil
}

func (c *Client) DeleteSession(ctx context.Context, key string) error {
	url := fmt.Sprintf("%s/api/v1/caches/session/%s", c.baseURL, key)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete session status %d", resp.StatusCode)
	}
	return nil
}

type sessionNotFoundError struct{}

func (sessionNotFoundError) Error() string { return "session not found" }

var ErrSessionNotFound = sessionNotFoundError{}

func IsSessionNotFound(err error) bool {
	_, ok := err.(sessionNotFoundError)
	return ok
}
