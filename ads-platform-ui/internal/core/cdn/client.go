package cdn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Category matches ads-platform-cdn GET /api/categories.
type Category struct {
	ID                             int     `json:"id"`
	Parent                         *int    `json:"parent"`
	Order                          int     `json:"order"`
	Title                          string  `json:"title"`
	Slug                           string  `json:"slug"`
	IsLeaf                         bool    `json:"isLeaf"`
	AdsAttrsJSONSchemaTemplateName *string `json:"adsAttrsJsonSchemaTemplateName"`
}

// City matches ads-platform-cdn GET /api/cities.
type City struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Parent      *int   `json:"parent"`
	Type        string `json:"type"`
	CitiesCount *int   `json:"cities_count,omitempty"`
}

// Client calls ads-platform-cdn HTTP APIs.
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

func (c *Client) GetCategories(ctx context.Context) ([]Category, error) {
	var items []Category
	if err := c.getJSON(ctx, "/api/categories", &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Client) GetCities(ctx context.Context) ([]City, error) {
	var items []City
	if err := c.getJSON(ctx, "/api/cities", &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Client) getJSON(ctx context.Context, path string, dest any) error {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("cdn: build request %s: %w", path, err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("cdn: GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cdn: GET %s: status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return fmt.Errorf("cdn: read %s: %w", url, err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("cdn: parse %s: %w", url, err)
	}
	return nil
}
