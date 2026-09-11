package impl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNominatimCityCenterUsesCatalog(t *testing.T) {
	c := NewNominatimClient("http://127.0.0.1:1")
	c.minInterval = 0
	pt, err := c.CityCenter(context.Background(), "tehran", "تهران", "fa")
	if err != nil {
		t.Fatal(err)
	}
	if pt.Source != "catalog" || pt.Lat < 35 || pt.Lat > 36 {
		t.Fatalf("%+v", pt)
	}
}

func TestNominatimReverse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/reverse" {
			t.Fatalf("path %s", r.URL.Path)
		}
		if r.Header.Get("User-Agent") == "" {
			t.Fatal("missing User-Agent")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"display_name": "Vanak, Tehran",
			"address": map[string]any{
				"suburb":  "Vanak",
				"city":    "Tehran",
				"country": "Iran",
			},
		})
	}))
	defer srv.Close()

	c := NewNominatimClient(srv.URL)
	c.minInterval = 0
	c.http = srv.Client()
	out, err := c.Reverse(context.Background(), 35.76, 51.41, "fa")
	if err != nil {
		t.Fatal(err)
	}
	if out.Neighborhood != "Vanak" {
		t.Fatalf("%+v", out)
	}

	// Cached: server would fail if hit again after close, so hit twice first.
	out2, err := c.Reverse(context.Background(), 35.76, 51.41, "fa")
	if err != nil || out2.Neighborhood != "Vanak" {
		t.Fatalf("cache: %v %+v", err, out2)
	}
}

func TestNominatimSearchFallback(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	c := NewNominatimClient(srv.URL)
	c.minInterval = 0
	c.http = &http.Client{Timeout: time.Second}
	pt, err := c.CityCenter(context.Background(), "unknown-city", "Nowhere", "fa")
	if err != nil {
		t.Fatal(err)
	}
	if pt.Source != "catalog" {
		t.Fatalf("expected tehran fallback, got %+v", pt)
	}
}
