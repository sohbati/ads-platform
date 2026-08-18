//go:build integration

package by_slugs

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/categorycache"
)

func TestCategoriesBySlugsAPI(t *testing.T) {
	_, cacheURL := categorycache.SetupCategoryCacheStack(t)

	t.Run("returns_categories_for_slugs", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, categories, errResp, err := categorycache.GetCategoriesBySlugs(ctx, cacheURL, "digital", "cars")
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("status=%d error=%+v", status, errResp)
		}
		if len(categories) != 2 {
			t.Fatalf("expected 2 categories, got %d: %+v", len(categories), categories)
		}

		bySlug := map[string]categorycache.Category{}
		for _, c := range categories {
			bySlug[c.Slug] = c
		}

		digital, ok := bySlug["digital"]
		if !ok {
			t.Fatal("missing digital")
		}
		if digital.ID != 3 {
			t.Fatalf("digital id: got %d", digital.ID)
		}
		if digital.Path != "3" {
			t.Fatalf("digital path: got %q", digital.Path)
		}

		cars, ok := bySlug["cars"]
		if !ok {
			t.Fatal("missing cars")
		}
		if cars.ID != 13 {
			t.Fatalf("cars id: got %d", cars.ID)
		}
		if cars.Path != "2,13" {
			t.Fatalf("cars path: got %q", cars.Path)
		}
	})

	t.Run("empty_slugs_returns_400", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, _, errResp, err := categorycache.GetCategoriesBySlugs(ctx, cacheURL)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d error=%+v", status, errResp)
		}
		if errResp.Error != "CATEGORY_SLUGS_EMPTY" {
			t.Fatalf("expected CATEGORY_SLUGS_EMPTY, got %q", errResp.Error)
		}
	})
}
