//go:build integration

package ads_api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"integration-test/internal/adsapi"

	"github.com/jackc/pgx/v5"
)

func TestGetPublicResidentialSaleReturnsAttrs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	_, ad, errResp, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, fullResidentialAttrs()))
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if ad.ID == 0 {
		t.Fatalf("create failed: %+v", errResp)
	}

	status, pub, errResp, err := adsapi.GetPublic(ctx, backURL, ad.ID)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if pub.CategoryID != residentialSaleCategoryID {
		t.Fatalf("category_id=%d", pub.CategoryID)
	}
	attrs := mustAttrs(t, pub.Attrs)
	if attrs["orientation"] != "north" || attrNum(t, attrs, "bathrooms") != 6 {
		t.Fatalf("public attrs=%s", pub.Attrs)
	}
	raw, _ := json.Marshal(pub)
	if strings.Contains(string(raw), "09121110001") || strings.Contains(string(raw), "9121110001") {
		t.Fatalf("full phone leaked: %s", raw)
	}
	if pub.MapLat == nil || pub.MapLng == nil {
		t.Fatal("expected approximate map coords")
	}
	if *pub.MapLat == 35.6892 && *pub.MapLng == 51.3890 {
		t.Fatal("public map coords must not be the exact pin")
	}
}

func TestGetPublicInactiveResidentialSaleNotFound(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	_, ad, _, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, nil))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	conn, err := pgx.Connect(ctx, pgHostDSN)
	if err != nil {
		t.Fatalf("postgres: %v", err)
	}
	defer conn.Close(ctx)
	if _, err := conn.Exec(ctx, `UPDATE ads_platform_schema.ads SET status='deleted' WHERE id=$1`, ad.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	status, _, errResp, err := adsapi.GetPublic(ctx, backURL, ad.ID)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusNotFound || errResp.Error != "AD_NOT_FOUND" {
		t.Fatalf("status=%d error=%+v, want 404 AD_NOT_FOUND", status, errResp)
	}
}

func TestSearchResidentialSaleCategory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	payload := residentialSalePayload(userID, map[string]any{"area_m2": 95, "rooms": 2})
	token := fmtToken()
	payload["title"] = "IT-RS-" + token
	_, ad, errResp, err := adsapi.PostJSON(ctx, backURL, payload)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if ad.ID == 0 {
		t.Fatalf("create failed: %+v", errResp)
	}

	status, found, errResp, err := adsapi.Search(ctx, backURL, "tehran", "residential-sale", "IT-RS-"+token)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if found.Pagination.Total < 1 {
		t.Fatalf("expected hit in residential-sale, got %+v", found)
	}
	hit := false
	for _, item := range found.Ads {
		if item.ID == ad.ID {
			hit = true
			break
		}
	}
	if !hit {
		t.Fatalf("ad %d not in results %+v", ad.ID, found.Ads)
	}

	status, other, errResp, err := adsapi.Search(ctx, backURL, "tehran", "laptop", "IT-RS-"+token)
	if err != nil {
		t.Fatalf("laptop search: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("laptop status=%d error=%+v", status, errResp)
	}
	for _, item := range other.Ads {
		if item.ID == ad.ID {
			t.Fatalf("residential ad leaked into laptop search")
		}
	}
}

func fmtToken() string {
	return time.Now().UTC().Format("150405.000000000")
}
