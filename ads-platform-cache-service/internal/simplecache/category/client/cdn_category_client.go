package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cache-service/internal/simplecache/category/model"
)

type CDNCategoryClient interface {
	ListCategories(ctx context.Context) ([]model.Category, error)
}

type httpCDNCategoryClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCDNCategoryClient(baseURL string) CDNCategoryClient {
	return &httpCDNCategoryClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *httpCDNCategoryClient) ListCategories(ctx context.Context) ([]model.Category, error) {
	url := c.baseURL + "/api/categories"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cdn categories status %d: %s", resp.StatusCode, string(body))
	}

	var categories []model.Category
	if err := json.Unmarshal(body, &categories); err != nil {
		return nil, err
	}
	return categories, nil
}
