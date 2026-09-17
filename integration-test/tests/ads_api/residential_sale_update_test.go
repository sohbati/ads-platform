//go:build integration

package ads_api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/adsapi"

	"github.com/jackc/pgx/v5"
)

func TestUpdateResidentialSaleAttrs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	_, ad, _, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, map[string]any{
		"area_m2": 80, "rooms": 1, "orientation": "north", "bathrooms": 1,
	}))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	payload := residentialSalePayload(userID, map[string]any{
		"area_m2": 95, "rooms": 2, "orientation": "south", "bathrooms": 2, "facade": "glass",
	})
	status, updated, errResp, err := adsapi.PutJSON(ctx, backURL, userID, ad.ID, payload)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if status != http.StatusOK {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	attrs := mustAttrs(t, updated.Attrs)
	if attrNum(t, attrs, "area_m2") != 95 || attrs["orientation"] != "south" || attrs["facade"] != "glass" {
		t.Fatalf("updated attrs=%s", updated.Attrs)
	}

	_, pub, _, err := adsapi.GetPublic(ctx, backURL, ad.ID)
	if err != nil {
		t.Fatalf("public: %v", err)
	}
	pubAttrs := mustAttrs(t, pub.Attrs)
	if pubAttrs["orientation"] != "south" || attrNum(t, pubAttrs, "bathrooms") != 2 {
		t.Fatalf("public attrs=%s", pub.Attrs)
	}
}

func TestUpdateResidentialSaleRejectsInvalidEnum(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	_, ad, _, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, map[string]any{
		"area_m2": 80, "rooms": 1, "deed_type": "single",
	}))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	payload := residentialSalePayload(userID, map[string]any{
		"area_m2": 80, "rooms": 1, "deed_type": "not-a-deed",
	})
	status, _, errResp, err := adsapi.PutJSON(ctx, backURL, userID, ad.ID, payload)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if status != http.StatusBadRequest || errResp.Error != "AD_INVALID_ATTRS" {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}

	conn, err := pgx.Connect(ctx, pgHostDSN)
	if err != nil {
		t.Fatalf("postgres: %v", err)
	}
	defer conn.Close(ctx)
	var stored []byte
	if err := conn.QueryRow(ctx, `SELECT attrs FROM ads_platform_schema.ads WHERE id=$1`, ad.ID).Scan(&stored); err != nil {
		t.Fatalf("select: %v", err)
	}
	attrs := mustAttrs(t, stored)
	if attrs["deed_type"] != "single" {
		t.Fatalf("attrs mutated after failed update: %s", stored)
	}
}

func TestUpdateResidentialSaleOtherUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	ownerID := insertUser(t, ctx)
	otherID := insertUser(t, ctx)
	_, ad, _, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(ownerID, nil))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	payload := residentialSalePayload(otherID, nil)
	payload["title"] = "Hijack"
	status, _, errResp, err := adsapi.PutJSON(ctx, backURL, otherID, ad.ID, payload)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if status != http.StatusNotFound || errResp.Error != "AD_NOT_FOUND" {
		t.Fatalf("status=%d error=%+v, want 404 AD_NOT_FOUND", status, errResp)
	}
}
