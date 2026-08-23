package categorycache

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

type Category struct {
	ID            int    `json:"id"`
	Parent        *int   `json:"parent"`
	Order         int    `json:"order"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Path          string `json:"path,omitempty"`
	IsLeaf        bool   `json:"isLeaf"`
	DescendantIDs []int  `json:"descendant_ids,omitempty"`
}

type ErrorResponse struct {
	Error      string   `json:"error"`
	StatusCode int      `json:"statusCode"`
	Params     []string `json:"params"`
}

func GetCategoriesByIDs(ctx context.Context, cacheURL string, ids ...int) (int, []Category, ErrorResponse, error) {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.Itoa(id))
	}
	u := cacheURL + "/api/v1/caches/categories/by-ids?ids=" + url.QueryEscape(strings.Join(parts, ","))
	return getCategories(ctx, u)
}

func GetCategoriesBySlugs(ctx context.Context, cacheURL string, slugs ...string) (int, []Category, ErrorResponse, error) {
	u := cacheURL + "/api/v1/caches/categories/by-slugs?slugs=" + url.QueryEscape(strings.Join(slugs, ","))
	return getCategories(ctx, u)
}

func GetCategoriesBySlugsWithDescendants(ctx context.Context, cacheURL string, slugs ...string) (int, []Category, ErrorResponse, error) {
	u := cacheURL + "/api/v1/caches/categories/by-slugs?include_descendants=true&slugs=" + url.QueryEscape(strings.Join(slugs, ","))
	return getCategories(ctx, u)
}

func getCategories(ctx context.Context, u string) (int, []Category, ErrorResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, nil, ErrorResponse{}, err
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, ErrorResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, ErrorResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		var failure ErrorResponse
		_ = json.Unmarshal(body, &failure)
		return resp.StatusCode, nil, failure, nil
	}

	var categories []Category
	if err := json.Unmarshal(body, &categories); err != nil {
		return resp.StatusCode, nil, ErrorResponse{}, fmt.Errorf("decode categories: %w; body=%s", err, string(body))
	}
	return resp.StatusCode, categories, ErrorResponse{}, nil
}
