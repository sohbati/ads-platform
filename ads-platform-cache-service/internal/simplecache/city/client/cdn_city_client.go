package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cache-service/internal/simplecache/city/model"
)

type CDNCityClient interface {
	ListCities(ctx context.Context) ([]model.City, error)
}

type httpCDNCityClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCDNCityClient(baseURL string) CDNCityClient {
	return &httpCDNCityClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *httpCDNCityClient) ListCities(ctx context.Context) ([]model.City, error) {
	url := c.baseURL + "/api/cities"
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
		return nil, fmt.Errorf("cdn cities status %d: %s", resp.StatusCode, string(body))
	}

	var cities []model.City
	if err := json.Unmarshal(body, &cities); err != nil {
		return nil, err
	}
	return cities, nil
}
