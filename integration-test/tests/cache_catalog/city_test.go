//go:build integration

package cache_catalog

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/citycache"
)

func TestCitiesByIDsAPI(t *testing.T) {
	t.Run("returns_cities_for_ids", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, cities, errResp, err := citycache.GetCitiesByIDs(ctx, cacheURL, 1, 869)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("status=%d error=%+v", status, errResp)
		}
		if len(cities) != 2 {
			t.Fatalf("expected 2 cities, got %d: %+v", len(cities), cities)
		}

		byID := map[int]citycache.City{}
		for _, c := range cities {
			byID[c.ID] = c
		}

		tehran, ok := byID[1]
		if !ok {
			t.Fatal("missing tehran id=1")
		}
		if tehran.Slug != "tehran" {
			t.Fatalf("tehran slug: got %q", tehran.Slug)
		}
		if tehran.Path == "" {
			t.Fatal("expected tehran path")
		}

		abyek, ok := byID[869]
		if !ok {
			t.Fatal("missing abyek id=869")
		}
		if abyek.Slug != "abyek" {
			t.Fatalf("abyek slug: got %q", abyek.Slug)
		}
	})

	t.Run("empty_ids_returns_400", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, _, errResp, err := citycache.GetCitiesByIDs(ctx, cacheURL)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d error=%+v", status, errResp)
		}
		if errResp.Error != "CITY_IDS_EMPTY" {
			t.Fatalf("expected CITY_IDS_EMPTY, got %q", errResp.Error)
		}
	})
}

func TestCitiesBySlugsAPI(t *testing.T) {
	t.Run("returns_cities_for_slugs", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, cities, errResp, err := citycache.GetCitiesBySlugs(ctx, cacheURL, "tehran", "abyek")
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("status=%d error=%+v", status, errResp)
		}
		if len(cities) != 2 {
			t.Fatalf("expected 2 cities, got %d: %+v", len(cities), cities)
		}

		bySlug := map[string]citycache.City{}
		for _, c := range cities {
			bySlug[c.Slug] = c
		}

		tehran, ok := bySlug["tehran"]
		if !ok {
			t.Fatal("missing tehran")
		}
		if tehran.ID != 1 {
			t.Fatalf("tehran id: got %d", tehran.ID)
		}

		abyek, ok := bySlug["abyek"]
		if !ok {
			t.Fatal("missing abyek")
		}
		if abyek.ID != 869 {
			t.Fatalf("abyek id: got %d", abyek.ID)
		}
	})

	t.Run("empty_slugs_returns_400", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, _, errResp, err := citycache.GetCitiesBySlugs(ctx, cacheURL)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d error=%+v", status, errResp)
		}
		if errResp.Error != "CITY_SLUGS_EMPTY" {
			t.Fatalf("expected CITY_SLUGS_EMPTY, got %q", errResp.Error)
		}
	})
}
