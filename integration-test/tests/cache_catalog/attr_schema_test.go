//go:build integration

package cache_catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/attrschemacache"
)

func TestAttrSchemasByNamesAPI(t *testing.T) {
	t.Run("returns_schemas_for_names", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, items, errResp, err := attrschemacache.GetAttrSchemasByNames(ctx, cacheURL, "cars", "laptop")
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusOK {
			t.Fatalf("status=%d error=%+v", status, errResp)
		}
		if len(items) != 2 {
			t.Fatalf("expected 2 schemas, got %d: %+v", len(items), items)
		}

		byName := map[string]attrschemacache.AttrSchema{}
		for _, item := range items {
			byName[item.Name] = item
		}

		cars, ok := byName["cars"]
		if !ok {
			t.Fatal("missing cars")
		}
		if cars.Title != "خودرو" {
			t.Fatalf("cars title: got %q", cars.Title)
		}
		var schema map[string]any
		if err := json.Unmarshal(cars.JSONSchema, &schema); err != nil {
			t.Fatalf("cars jsonSchema: %v", err)
		}
		props, _ := schema["properties"].(map[string]any)
		if _, ok := props["gearbox"]; !ok {
			t.Fatalf("cars schema missing gearbox: %s", cars.JSONSchema)
		}

		laptop, ok := byName["laptop"]
		if !ok {
			t.Fatal("missing laptop")
		}
		if laptop.Title != "لپ‌تاپ" {
			t.Fatalf("laptop title: got %q", laptop.Title)
		}
	})

	t.Run("empty_names_returns_400", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		status, _, errResp, err := attrschemacache.GetAttrSchemasByNames(ctx, cacheURL)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		if status != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d error=%+v", status, errResp)
		}
		if errResp.Error != "ATTR_SCHEMA_NAMES_EMPTY" {
			t.Fatalf("expected ATTR_SCHEMA_NAMES_EMPTY, got %q", errResp.Error)
		}
	})
}
