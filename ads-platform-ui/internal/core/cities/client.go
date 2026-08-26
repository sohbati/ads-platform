package cities

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const citiesAPIPath = "/api/cities"

// FetchFromCDN downloads and indexes the cities list from ads-platform-cdn.
func FetchFromCDN(ctx context.Context, cdnBaseURL string, client *http.Client) (*Catalog, error) {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	url := strings.TrimRight(strings.TrimSpace(cdnBaseURL), "/") + citiesAPIPath
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("cities: build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cities: GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cities: GET %s: status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 32<<20))
	if err != nil {
		return nil, fmt.Errorf("cities: read response: %w", err)
	}

	return parseRecords(data, url)
}

func parseRecords(data []byte, source string) (*Catalog, error) {
	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("cities: parse %s: %w", source, err)
	}
	c := &Catalog{
		bySlug:   make(map[string]Record, len(records)),
		byID:     make(map[int]Record, len(records)),
		children: make(map[int][]Record),
		all:      records,
	}
	for _, r := range records {
		c.byID[r.ID] = r
		if r.Slug != "" {
			c.bySlug[normalizeSlug(r.Slug)] = r
		}
		if r.Parent != nil {
			c.children[*r.Parent] = append(c.children[*r.Parent], r)
		}
	}
	if len(c.bySlug) == 0 {
		return nil, fmt.Errorf("cities: no entries in %s", source)
	}
	return c, nil
}
