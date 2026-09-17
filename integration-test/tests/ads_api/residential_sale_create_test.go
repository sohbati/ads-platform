//go:build integration

package ads_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"integration-test/internal/adsapi"

	"github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	residentialSaleCategoryID = 11
	realEstateParentID        = 1
)

func residentialSalePayload(userID int64, attrs map[string]any) map[string]any {
	if attrs == nil {
		attrs = map[string]any{"area_m2": 90, "rooms": 2}
	}
	return map[string]any{
		"user_id":      userID,
		"category_id":  residentialSaleCategoryID,
		"city_id":      tehranCityID,
		"title":        "آپارتمان ۹۰ متری سعادت‌آباد",
		"description":  "فروش آپارتمان نوساز.",
		"latitude":     35.6892,
		"longitude":    51.3890,
		"neighborhood": "Saadat Abad",
		"price_amount": 12_000_000_000,
		"price_type":   "fixed",
		"currency":     "IRR",
		"attrs":        attrs,
		"contact":      map[string]any{"phone": "09121110001", "chat_enabled": true},
	}
}

func fullResidentialAttrs() map[string]any {
	return map[string]any{
		"area_m2":         90,
		"rooms":           2,
		"parking":         true,
		"elevator":        false,
		"floor":           3,
		"balcony":         true,
		"year_built":      1401,
		"has_storage":     true,
		"storage_m2":      8.5,
		"floors":          5,
		"units":           10,
		"orientation":     "north",
		"flooring":        "ceramic",
		"carpeted":        false,
		"toilets_iranian": 1,
		"toilets_western": 0,
		"bathrooms":       6,
		"hot_water":       "package",
		"deed_type":       "booklet",
		"facade":          "stone",
		"property_price":  999,
		"brand":           "lenovo",
	}
}

func mustAttrs(t *testing.T, raw json.RawMessage) map[string]any {
	t.Helper()
	attrs, err := adsapi.ParseAttrs(raw)
	if err != nil {
		t.Fatalf("parse attrs: %v (%s)", err, raw)
	}
	return attrs
}

func attrNum(t *testing.T, attrs map[string]any, key string) float64 {
	t.Helper()
	n, ok := attrs[key].(float64)
	if !ok {
		t.Fatalf("attr %s = %#v", key, attrs[key])
	}
	return n
}

func TestCreateResidentialSaleJSON(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	status, ad, errResp, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, fullResidentialAttrs()))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if ad.Status != "active" || ad.CategoryID != residentialSaleCategoryID || ad.CityID != tehranCityID {
		t.Fatalf("ad=%+v", ad)
	}

	attrs := mustAttrs(t, ad.Attrs)
	if attrNum(t, attrs, "area_m2") != 90 || attrNum(t, attrs, "rooms") != 2 {
		t.Fatalf("required attrs=%s", ad.Attrs)
	}
	if attrs["orientation"] != "north" || attrs["deed_type"] != "booklet" || attrs["hot_water"] != "package" {
		t.Fatalf("enums=%s", ad.Attrs)
	}
	if attrs["parking"] != true || attrs["elevator"] != false || attrs["carpeted"] != false {
		t.Fatalf("booleans=%s", ad.Attrs)
	}
	if attrNum(t, attrs, "storage_m2") != 8.5 || attrNum(t, attrs, "toilets_western") != 0 || attrNum(t, attrs, "bathrooms") != 6 {
		t.Fatalf("numbers=%s", ad.Attrs)
	}
	if _, ok := attrs["property_price"]; ok {
		t.Fatalf("unknown key stored: %s", ad.Attrs)
	}
	if _, ok := attrs["brand"]; ok {
		t.Fatalf("laptop attr stored: %s", ad.Attrs)
	}

	var stored json.RawMessage
	conn, err := pgx.Connect(ctx, pgHostDSN)
	if err != nil {
		t.Fatalf("postgres: %v", err)
	}
	defer conn.Close(ctx)
	if err := conn.QueryRow(ctx, `SELECT attrs FROM ads_platform_schema.ads WHERE id=$1`, ad.ID).Scan(&stored); err != nil {
		t.Fatalf("select attrs: %v", err)
	}
	pgAttrs := mustAttrs(t, stored)
	if pgAttrs["orientation"] != "north" || attrNum(t, pgAttrs, "bathrooms") != 6 {
		t.Fatalf("postgres attrs=%s", stored)
	}
}

func TestCreateResidentialSaleRequiredOnly(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	status, ad, errResp, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, map[string]any{
		"area_m2": 80, "rooms": 1,
	}))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	attrs := mustAttrs(t, ad.Attrs)
	if len(attrs) != 2 {
		t.Fatalf("attrs=%s, want only area_m2 and rooms", ad.Attrs)
	}
}

func TestCreateResidentialSaleBooleans(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	status, ad, errResp, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, map[string]any{
		"area_m2": 70, "rooms": 1,
		"parking": false, "elevator": false, "balcony": false, "has_storage": false, "carpeted": false,
	}))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	attrs := mustAttrs(t, ad.Attrs)
	for _, key := range []string{"parking", "elevator", "balcony", "has_storage", "carpeted"} {
		if attrs[key] != false {
			t.Fatalf("%s=%#v, want false in %s", key, attrs[key], ad.Attrs)
		}
	}
}

func TestCreateResidentialSaleBathroomSelects(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	status, ad, errResp, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, map[string]any{
		"area_m2": 100, "rooms": 3,
		"toilets_iranian": 0, "toilets_western": 1, "bathrooms": 6,
	}))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	attrs := mustAttrs(t, ad.Attrs)
	if attrNum(t, attrs, "toilets_iranian") != 0 || attrNum(t, attrs, "toilets_western") != 1 || attrNum(t, attrs, "bathrooms") != 6 {
		t.Fatalf("attrs=%s", ad.Attrs)
	}
}

func TestCreateResidentialSaleDropsUnknownAttrs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	status, ad, errResp, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, map[string]any{
		"area_m2": 90, "rooms": 2, "foo": "bar", "property_price": 1,
	}))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	attrs := mustAttrs(t, ad.Attrs)
	if _, ok := attrs["foo"]; ok {
		t.Fatalf("foo stored: %s", ad.Attrs)
	}
	if _, ok := attrs["property_price"]; ok {
		t.Fatalf("property_price stored: %s", ad.Attrs)
	}
}

func TestCreateResidentialSaleWithPictures(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	status, ad, errResp, err := adsapi.PostMultipart(ctx, backURL, residentialSalePayload(userID, nil), map[string][]byte{
		"cover.jpg": adsapi.JPEG1x1,
	})
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	items, err := adsapi.ParseMedia(ad.Media)
	if err != nil || len(items) != 1 {
		t.Fatalf("media=%s err=%v", ad.Media, err)
	}
	mc, err := minio.New(minioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(minioUser, minioPass, ""),
		Secure: false,
	})
	if err != nil {
		t.Fatalf("minio: %v", err)
	}
	if _, err := mc.StatObject(ctx, minioBucket, items[0].ObjectKey, minio.StatObjectOptions{}); err != nil {
		t.Fatalf("stat object: %v", err)
	}
}

func TestCreateResidentialSaleValidationErrors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	userID := insertUser(t, ctx)

	cases := []struct {
		name     string
		mutate   func(map[string]any)
		wantCode string
	}{
		{"missing area_m2", func(p map[string]any) {
			p["attrs"] = map[string]any{"rooms": 2}
		}, "AD_INVALID_ATTRS"},
		{"missing rooms", func(p map[string]any) {
			p["attrs"] = map[string]any{"area_m2": 90}
		}, "AD_INVALID_ATTRS"},
		{"area_m2 zero", func(p map[string]any) {
			p["attrs"] = map[string]any{"area_m2": 0, "rooms": 2}
		}, "AD_INVALID_ATTRS"},
		{"rooms above max", func(p map[string]any) {
			p["attrs"] = map[string]any{"area_m2": 90, "rooms": 51}
		}, "AD_INVALID_ATTRS"},
		{"floors below min", func(p map[string]any) {
			p["attrs"] = map[string]any{"area_m2": 90, "rooms": 2, "floors": 0}
		}, "AD_INVALID_ATTRS"},
		{"invalid orientation", func(p map[string]any) {
			p["attrs"] = map[string]any{"area_m2": 90, "rooms": 2, "orientation": "west"}
		}, "AD_INVALID_ATTRS"},
		{"invalid deed", func(p map[string]any) {
			p["attrs"] = map[string]any{"area_m2": 90, "rooms": 2, "deed_type": "foo"}
		}, "AD_INVALID_ATTRS"},
		{"bathrooms 7", func(p map[string]any) {
			p["attrs"] = map[string]any{"area_m2": 90, "rooms": 2, "bathrooms": 7}
		}, "AD_INVALID_ATTRS"},
		{"has_storage string", func(p map[string]any) {
			p["attrs"] = map[string]any{"area_m2": 90, "rooms": 2, "has_storage": "yes"}
		}, "AD_INVALID_ATTRS"},
		{"attrs array", func(p map[string]any) {
			p["attrs"] = []any{"nope"}
		}, "AD_INVALID_ATTRS"},
		{"parent category", func(p map[string]any) {
			p["category_id"] = realEstateParentID
		}, "AD_CATEGORY_NOT_LEAF"},
		{"unknown category", func(p map[string]any) {
			p["category_id"] = 999999
		}, "AD_INVALID_CATEGORY"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := residentialSalePayload(userID, map[string]any{"area_m2": 90, "rooms": 2})
			tc.mutate(payload)
			status, _, errResp, err := adsapi.PostJSON(ctx, backURL, payload)
			if err != nil {
				t.Fatalf("request: %v", err)
			}
			if status != http.StatusBadRequest || errResp.Error != tc.wantCode {
				t.Fatalf("status=%d error=%+v, want 400 %s", status, errResp, tc.wantCode)
			}
		})
	}
}

func TestCreateResidentialSaleRejectsNonObjectAttrs(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	userID := insertUser(t, ctx)

	body := []byte(fmt.Sprintf(`{
		"user_id": %d,
		"category_id": %d,
		"city_id": %d,
		"title": "x",
		"description": "y",
		"attrs": "{"
	}`, userID, residentialSaleCategoryID, tehranCityID))
	status, _, errResp, err := adsapi.PostRawJSON(ctx, backURL, body)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusBadRequest || errResp.Error != "AD_INVALID_ATTRS" {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
}

func TestCreateResidentialSaleLaptopAttrsNeedRequired(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	userID := insertUser(t, ctx)

	status, _, errResp, err := adsapi.PostJSON(ctx, backURL, residentialSalePayload(userID, map[string]any{
		"brand": "lenovo", "ram_gb": 16,
	}))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusBadRequest || errResp.Error != "AD_INVALID_ATTRS" {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
}
