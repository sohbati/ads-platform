//go:build integration

package by_slugs

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/citycache"
)

func TestCitiesBySlugsAPI(t *testing.T) {
	_, cacheURL := citycache.SetupCityCacheStack(t)

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
		if tehran.Path == "" {
			t.Fatal("expected tehran path")
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
