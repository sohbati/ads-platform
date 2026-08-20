package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Category mirrors the cache-service category response.
type Category struct {
	ID            int    `json:"id"`
	Slug          string `json:"slug"`
	Title         string `json:"title"`
	Path          string `json:"path,omitempty"`
	DescendantIDs []int  `json:"descendant_ids,omitempty"`
}

// City mirrors the cache-service city response.
type City struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// CatalogClient resolves categories and cities via ads-platform-cache-service.
type CatalogClient interface {
	CategoriesBySlugs(ctx context.Context, slugs []string, includeDescendants bool) ([]Category, error)
	CitiesBySlugs(ctx context.Context, slugs []string) ([]City, error)
	CitiesByIDs(ctx context.Context, ids []int) ([]City, error)
}

type catalogClient struct {
	baseURL string
	http    *http.Client
}

func NewCatalogClient(baseURL string, httpClient *http.Client) CatalogClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &catalogClient{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		http:    httpClient,
	}
}

func (c *catalogClient) CategoriesBySlugs(ctx context.Context, slugs []string, includeDescendants bool) ([]Category, error) {
	q := url.Values{"slugs": {strings.Join(slugs, ",")}}
	if includeDescendants {
		q.Set("include_descendants", "true")
	}

	var out []Category
	if err := c.getJSON(ctx, "/api/v1/caches/categories/by-slugs", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogClient) CitiesBySlugs(ctx context.Context, slugs []string) ([]City, error) {
	q := url.Values{"slugs": {strings.Join(slugs, ",")}}

	var out []City
	if err := c.getJSON(ctx, "/api/v1/caches/cities/by-slugs", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogClient) CitiesByIDs(ctx context.Context, ids []int) ([]City, error) {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.Itoa(id))
	}
	q := url.Values{"ids": {strings.Join(parts, ",")}}

	var out []City
	if err := c.getJSON(ctx, "/api/v1/caches/cities/by-ids", q, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *catalogClient) getJSON(ctx context.Context, path string, query url.Values, out any) error {
	reqURL := c.baseURL + path + "?" + query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("cache-service: build request %s: %w", reqURL, err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("cache-service: GET %s: %w", reqURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("cache-service: read %s: %w", reqURL, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cache-service: GET %s: status %d: %s", reqURL, resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("cache-service: parse %s: %w", reqURL, err)
	}
	return nil
}
