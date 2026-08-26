package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"ads-platform-ui/internal/core/bff"
)

// SearchResponse mirrors ads-platform-back GET /api/v1/q/:place/:category.
type SearchResponse struct {
	Place         string     `json:"place"`
	Category      string     `json:"category"`
	CategoryTitle string     `json:"category_title"`
	Pagination    Pagination `json:"pagination"`
	Ads           []Ad       `json:"ads"`
}

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type Ad struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	PriceAmount  *int64  `json:"price_amount"`
	PriceType    string  `json:"price_type"`
	Currency     string  `json:"currency"`
	CityName     string  `json:"city_name"`
	Neighborhood string  `json:"neighborhood"`
	Thumbnail    string  `json:"thumbnail"`
	HasPhoto     bool    `json:"has_photo"`
	PublishedAt  *string `json:"published_at"`
}

// ErrNotFound marks 4xx responses (unknown place/category, bad params) so the
// page can show an empty state instead of a service failure.
var ErrNotFound = fmt.Errorf("search: no results for request")

// SearchClient queries the ads search API through ads-bff.
type SearchClient struct {
	bff *bff.Client
}

func NewSearchClient(bffClient *bff.Client) *SearchClient {
	return &SearchClient{bff: bffClient}
}

func (c *SearchClient) Search(ctx context.Context, place, category, query string, page int, citiesCSV string) (*SearchResponse, error) {
	q := url.Values{}
	if query != "" {
		q.Set("q", query)
	}
	if page > 1 {
		q.Set("page", strconv.Itoa(page))
	}
	if citiesCSV != "" {
		q.Set("cities", citiesCSV)
	}

	path := "/api/v1/q/" + url.PathEscape(place) + "/" + url.PathEscape(category)
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}

	result, err := c.bff.Get(ctx, path, "")
	if err != nil {
		return nil, err
	}
	if result.StatusCode >= 400 && result.StatusCode < 500 {
		return nil, ErrNotFound
	}
	if result.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search: %s: status %d", path, result.StatusCode)
	}

	var out SearchResponse
	if err := json.Unmarshal(result.Body, &out); err != nil {
		return nil, fmt.Errorf("search: parse %s: %w", path, err)
	}
	return &out, nil
}
