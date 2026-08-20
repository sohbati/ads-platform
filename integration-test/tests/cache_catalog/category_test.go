//go:build integration

package cache_catalog

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/categorycache"
)

func TestCategoriesByIDsAPI(t *testing.T) {
	t.Run("returns_categories_for_ids", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

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

func TestCategoriesBySlugsAPI(t *testing.T) {
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

		cars, ok := bySlug["cars"]
		if !ok {
			t.Fatal("missing cars")
		}
		if cars.ID != 13 {
			t.Fatalf("cars id: got %d", cars.ID)
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

	t.Run("include_descendants_returns_subtree_ids", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, categories, errResp, err := categorycache.GetCategoriesBySlugsWithDescendants(ctx, cacheURL, "digital")
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("status=%d error=%+v", status, errResp)
		}
		if len(categories) != 1 {
			t.Fatalf("expected 1 category, got %d: %+v", len(categories), categories)
		}

		digital := categories[0]
		if digital.ID != 3 {
			t.Fatalf("digital id: got %d", digital.ID)
		}
		if len(digital.DescendantIDs) < 2 {
			t.Fatalf("expected descendants for digital, got %v", digital.DescendantIDs)
		}

		want := map[int]bool{3: false, 15: false} // self + mobile-tablet
		for _, id := range digital.DescendantIDs {
			if _, ok := want[id]; ok {
				want[id] = true
			}
		}
		for id, found := range want {
			if !found {
				t.Fatalf("expected id %d in descendants: %v", id, digital.DescendantIDs)
			}
		}
	})
}
