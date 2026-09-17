//go:build integration

package cache_catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/attrschemacache"
	"integration-test/internal/categorycache"
)

func TestResidentialSaleCategoryBySlug(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	status, categories, errResp, err := categorycache.GetCategoriesBySlugs(ctx, cacheURL, "residential-sale")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if len(categories) != 1 {
		t.Fatalf("expected 1 category, got %d: %+v", len(categories), categories)
	}
	cat := categories[0]
	if cat.ID != 11 || cat.Slug != "residential-sale" || !cat.IsLeaf {
		t.Fatalf("category: %+v", cat)
	}
	if cat.Parent == nil || *cat.Parent != 1 {
		t.Fatalf("parent=%v, want real-estate id 1", cat.Parent)
	}
	if cat.AdsAttrsJSONSchemaTemplateName == nil || *cat.AdsAttrsJSONSchemaTemplateName != "residential-sale" {
		t.Fatalf("template=%v", cat.AdsAttrsJSONSchemaTemplateName)
	}
	if cat.Path != "1,11" {
		t.Fatalf("path=%q", cat.Path)
	}
}

func TestResidentialSaleCategoryByID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	status, categories, errResp, err := categorycache.GetCategoriesByIDs(ctx, cacheURL, 11)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if len(categories) != 1 || categories[0].Slug != "residential-sale" {
		t.Fatalf("categories=%+v", categories)
	}
}

func TestRealEstateDescendantsIncludeResidentialSale(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	status, categories, errResp, err := categorycache.GetCategoriesBySlugsWithDescendants(ctx, cacheURL, "real-estate")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if len(categories) != 1 || categories[0].ID != 1 {
		t.Fatalf("categories=%+v", categories)
	}
	found := false
	for _, id := range categories[0].DescendantIDs {
		if id == 11 {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected id 11 in descendants: %v", categories[0].DescendantIDs)
	}
}

func TestResidentialSaleAttrSchemaShape(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	status, items, errResp, err := attrschemacache.GetAttrSchemasByNames(ctx, cacheURL, "residential-sale")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 schema, got %d", len(items))
	}
	schema := items[0]
	if schema.Name != "residential-sale" || schema.Title != "فروش مسکونی" {
		t.Fatalf("schema name/title: %+v", schema)
	}

	var doc map[string]any
	if err := json.Unmarshal(schema.JSONSchema, &doc); err != nil {
		t.Fatalf("jsonSchema: %v", err)
	}
	props, _ := doc["properties"].(map[string]any)
	for _, name := range []string{
		"area_m2", "rooms", "parking", "elevator", "balcony", "has_storage",
		"storage_m2", "orientation", "flooring", "toilets_iranian", "toilets_western",
		"bathrooms", "hot_water", "deed_type", "facade",
	} {
		if _, ok := props[name]; !ok {
			t.Fatalf("missing property %s in %s", name, schema.JSONSchema)
		}
	}
	required, _ := doc["required"].([]any)
	if !containsAny(required, "area_m2") || !containsAny(required, "rooms") {
		t.Fatalf("required=%v", required)
	}
	bathrooms, _ := props["bathrooms"].(map[string]any)
	enum, _ := bathrooms["enum"].([]any)
	if len(enum) != 7 {
		t.Fatalf("bathrooms enum=%v, want 0..6", enum)
	}
}

func containsAny(items []any, want string) bool {
	for _, item := range items {
		if s, ok := item.(string); ok && s == want {
			return true
		}
	}
	return false
}
