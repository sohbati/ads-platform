//go:build integration

package ads_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	tc "integration-test/internal/testcontainers"

	"github.com/jackc/pgx/v5"
)

func TestAdsTableInsertAndJSONBQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	net, err := tc.CreateNetwork(ctx)
	if err != nil {
		t.Fatalf("create network: %v", err)
	}
	t.Cleanup(func() {
		_ = net.Remove(context.Background())
	})

	pg, err := tc.StartPostgres(ctx, net)
	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}
	t.Cleanup(func() {
		_ = pg.Container.Terminate(context.Background())
	})

	conn, err := pgx.Connect(ctx, pg.HostDSN)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}
	defer conn.Close(ctx)

	var userID int64
	err = conn.QueryRow(ctx, `
		INSERT INTO ads_platform_schema."user" (name, mobile, national_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, "Test Seller", "09120000000", "0012345678").Scan(&userID)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	attrs := map[string]any{
		"brand":      "peugeot",
		"model":      "206",
		"year":       1399,
		"mileage_km": 87000,
	}
	media := []map[string]any{
		{"url": "https://cdn.example/a.jpg", "order": 1, "is_cover": true},
	}
	contact := map[string]any{"phone": "09120000000", "chat_enabled": true}
	location := map[string]any{"lat": 35.6892, "lng": 51.3890, "neighborhood": "vanak"}

	attrsJSON, _ := json.Marshal(attrs)
	mediaJSON, _ := json.Marshal(media)
	contactJSON, _ := json.Marshal(contact)
	locationJSON, _ := json.Marshal(location)

	var adID int64
	err = conn.QueryRow(ctx, `
		INSERT INTO ads_platform_schema.ads (
			user_id, category_id, city_id, title, description, status,
			price_amount, price_type, currency, attrs, media, contact, location, slug
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10::jsonb, $11::jsonb, $12::jsonb, $13::jsonb, $14
		)
		RETURNING id
	`,
		userID, 2, 1, "Peugeot 206", "Clean car, one owner", "active",
		450000000, "fixed", "IRR", attrsJSON, mediaJSON, contactJSON, locationJSON, "peugeot-206",
	).Scan(&adID)
	if err != nil {
		t.Fatalf("insert ad: %v", err)
	}
	if adID <= 0 {
		t.Fatalf("expected positive ad id, got %d", adID)
	}

	var (
		gotTitle  string
		gotBrand  string
		gotStatus string
		gotCover  bool
	)
	err = conn.QueryRow(ctx, `
		SELECT
			title,
			attrs->>'brand',
			status,
			(media->0->>'is_cover')::boolean
		FROM ads_platform_schema.ads
		WHERE id = $1
		  AND attrs @> '{"brand":"peugeot"}'::jsonb
	`, adID).Scan(&gotTitle, &gotBrand, &gotStatus, &gotCover)
	if err != nil {
		t.Fatalf("query ad jsonb: %v", err)
	}

	if gotTitle != "Peugeot 206" {
		t.Fatalf("title: got %q", gotTitle)
	}
	if gotBrand != "peugeot" {
		t.Fatalf("brand: got %q", gotBrand)
	}
	if gotStatus != "active" {
		t.Fatalf("status: got %q", gotStatus)
	}
	if !gotCover {
		t.Fatal("expected cover image")
	}

	var indexCount int
	err = conn.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM pg_indexes
		WHERE schemaname = 'ads_platform_schema'
		  AND tablename = 'ads'
		  AND indexname IN ('ads_attrs_gin_idx', 'ads_city_status_idx')
	`).Scan(&indexCount)
	if err != nil {
		t.Fatalf("query indexes: %v", err)
	}
	if indexCount != 2 {
		t.Fatalf("expected ads indexes, got count=%d", indexCount)
	}
}
