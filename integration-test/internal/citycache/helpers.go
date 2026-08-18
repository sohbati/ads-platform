package citycache

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

type City struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Parent      *int   `json:"parent"`
	Type        string `json:"type"`
	CitiesCount *int   `json:"cities_count,omitempty"`
	Path        string `json:"path,omitempty"`
}

type ErrorResponse struct {
	Error      string   `json:"error"`
	StatusCode int      `json:"statusCode"`
	Params     []string `json:"params"`
}

func GetCitiesByIDs(ctx context.Context, cacheURL string, ids ...int) (int, []City, ErrorResponse, error) {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.Itoa(id))
	}
	u := cacheURL + "/api/v1/caches/cities/by-ids?ids=" + url.QueryEscape(strings.Join(parts, ","))
	return getCities(ctx, u)
}

func GetCitiesBySlugs(ctx context.Context, cacheURL string, slugs ...string) (int, []City, ErrorResponse, error) {
	u := cacheURL + "/api/v1/caches/cities/by-slugs?slugs=" + url.QueryEscape(strings.Join(slugs, ","))
	return getCities(ctx, u)
}

func getCities(ctx context.Context, u string) (int, []City, ErrorResponse, error) {
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

	var cities []City
	if err := json.Unmarshal(body, &cities); err != nil {
		return resp.StatusCode, nil, ErrorResponse{}, fmt.Errorf("decode cities: %w; body=%s", err, string(body))
	}
	return resp.StatusCode, cities, ErrorResponse{}, nil
}
