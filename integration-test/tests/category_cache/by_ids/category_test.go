//go:build integration

package by_ids

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/categorycache"
)

func TestCategoriesByIDsAPI(t *testing.T) {
	_, cacheURL := categorycache.SetupCategoryCacheStack(t)

	t.Run("returns_categories_for_ids", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		// 3=digital, 13=cars — by-ids uses id→slug then slug→category
		status, categories, errResp, err := categorycache.GetCategoriesByIDs(ctx, cacheURL, 3, 13)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("status=%d error=%+v", status, errResp)
		}
		if len(categories) != 2 {
			t.Fatalf("expected 2 categories, got %d: %+v", len(categories), categories)
		}

		byID := map[int]categorycache.Category{}
		for _, c := range categories {
			byID[c.ID] = c
		}

		digital, ok := byID[3]
		if !ok {
			t.Fatal("missing digital id=3")
		}
		if digital.Slug != "digital" {
			t.Fatalf("digital slug: got %q", digital.Slug)
		}
		if digital.Path != "3" {
			t.Fatalf("digital path: got %q", digital.Path)
		}

		cars, ok := byID[13]
		if !ok {
			t.Fatal("missing cars id=13")
		}
		if cars.Slug != "cars" {
			t.Fatalf("cars slug: got %q", cars.Slug)
		}
		if cars.Path != "2,13" {
			t.Fatalf("cars path: got %q", cars.Path)
		}
	})

	t.Run("empty_ids_returns_400", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, _, errResp, err := categorycache.GetCategoriesByIDs(ctx, cacheURL)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d error=%+v", status, errResp)
		}
		if errResp.Error != "CATEGORY_IDS_EMPTY" {
			t.Fatalf("expected CATEGORY_IDS_EMPTY, got %q", errResp.Error)
		}
	})
}
