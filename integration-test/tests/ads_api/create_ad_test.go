//go:build integration

package ads_api

import (
	"bytes"
	"context"
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
	leafCategoryID   = 16 // laptop — no children in category.json
	parentCategoryID = 2  // vehicles
	tehranCityID     = 1
)

func insertUser(t *testing.T, ctx context.Context) int64 {
	t.Helper()
	conn, err := pgx.Connect(ctx, pgHostDSN)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}
	defer conn.Close(ctx)

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	var id int64
	err = conn.QueryRow(ctx, `
		INSERT INTO ads_platform_schema."user" (name, mobile, national_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, "Seller "+suffix, "09"+suffix[len(suffix)-9:], suffix).Scan(&id)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	return id
}

func validPayload(userID int64) map[string]any {
	return map[string]any{
		"user_id":      userID,
		"category_id":  leafCategoryID,
		"city_id":      tehranCityID,
		"title":        "Used laptop",
		"description":  "Lightly used\nwith charger.",
		"latitude":     35.6892,
		"longitude":    51.3890,
		"neighborhood": "Vanak",
		"price_amount": 12_000_000,
		"price_type":   "fixed",
		"currency":     "IRR",
		"attrs":        map[string]any{"brand": "lenovo", "ram_gb": 16},
		"contact":      map[string]any{"chat_enabled": true},
	}
}

func TestCreateAdJSON(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	status, ad, errResp, err := adsapi.PostJSON(ctx, backURL, validPayload(userID))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}
	if ad.ID <= 0 || ad.UserID != userID {
		t.Fatalf("ad identity: %+v", ad)
	}
	if ad.Status != "active" {
		t.Fatalf("status=%q, want active", ad.Status)
	}
	if ad.Title != "Used laptop" || ad.Description != "Lightly used\nwith charger." {
		t.Fatalf("title/description: %q / %q", ad.Title, ad.Description)
	}
	if ad.CategoryID != leafCategoryID || ad.CityID != tehranCityID {
		t.Fatalf("catalog ids: category=%d city=%d", ad.CategoryID, ad.CityID)
	}
	if !bytes.Contains(ad.Location, []byte(`"lat":35.6892`)) || !bytes.Contains(ad.Location, []byte("Vanak")) {
		t.Fatalf("location=%s", ad.Location)
	}
	if !bytes.Contains(ad.Attrs, []byte(`"brand":"lenovo"`)) {
		t.Fatalf("attrs=%s", ad.Attrs)
	}
	if string(ad.Media) != "[]" {
		t.Fatalf("media=%s, want []", ad.Media)
	}
}

func TestCreateAdWithPicturesUploadsToMinio(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	status, ad, errResp, err := adsapi.PostMultipart(ctx, backURL, validPayload(userID), map[string][]byte{
		"cover.jpg": adsapi.JPEG1x1,
	})
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusCreated {
		t.Fatalf("status=%d error=%+v", status, errResp)
	}

	items, err := adsapi.ParseMedia(ad.Media)
	if err != nil {
		t.Fatalf("parse media: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 media item, got %d (%s)", len(items), ad.Media)
	}
	item := items[0]
	if item.ObjectKey != fmt.Sprintf("ads/%d/%d_1.jpg", userID, userID) || !item.IsCover || item.ContentType != "image/jpeg" {
		t.Fatalf("media item: %+v", item)
	}

	mc, err := minio.New(minioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(minioUser, minioPass, ""),
		Secure: false,
	})
	if err != nil {
		t.Fatalf("minio client: %v", err)
	}
	info, err := mc.StatObject(ctx, minioBucket, item.ObjectKey, minio.StatObjectOptions{})
	if err != nil {
		t.Fatalf("stat minio object %s: %v", item.ObjectKey, err)
	}
	if info.Size != int64(len(adsapi.JPEG1x1)) {
		t.Fatalf("object size=%d, want %d", info.Size, len(adsapi.JPEG1x1))
	}

	conn, err := pgx.Connect(ctx, pgHostDSN)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}
	defer conn.Close(ctx)

	var (
		imageAdID int64
		statusRow string
		fileSize  int64
	)
	err = conn.QueryRow(ctx, `
		SELECT ad_id, status, file_size
		FROM ads_platform_schema.ad_images
		WHERE object_key = $1
	`, item.ObjectKey).Scan(&imageAdID, &statusRow, &fileSize)
	if err != nil {
		t.Fatalf("query ad_images: %v", err)
	}
	if imageAdID != ad.ID || statusRow != "uploaded" || fileSize != int64(len(adsapi.JPEG1x1)) {
		t.Fatalf("ad_images row: ad_id=%d status=%q size=%d", imageAdID, statusRow, fileSize)
	}
}

func TestCreateAdRejectsParentCategory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	payload := validPayload(userID)
	payload["category_id"] = parentCategoryID

	status, _, errResp, err := adsapi.PostJSON(ctx, backURL, payload)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusBadRequest || errResp.Error != "AD_CATEGORY_NOT_LEAF" {
		t.Fatalf("status=%d error=%+v, want 400 AD_CATEGORY_NOT_LEAF", status, errResp)
	}
}

func TestCreateAdTooManyPictures(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	userID := insertUser(t, ctx)
	status, _, errResp, err := adsapi.PostMultipart(ctx, backURL, validPayload(userID), map[string][]byte{
		"a.jpg": adsapi.JPEG1x1,
		"b.jpg": adsapi.JPEG1x1,
		"c.jpg": adsapi.JPEG1x1,
	})
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	if status != http.StatusBadRequest || errResp.Error != "AD_TOO_MANY_PICTURES" {
		t.Fatalf("status=%d error=%+v, want 400 AD_TOO_MANY_PICTURES", status, errResp)
	}
}

func TestCreateAdValidationErrors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	userID := insertUser(t, ctx)

	cases := []struct {
		name     string
		mutate   func(map[string]any)
		wantCode string
	}{
		{"empty title", func(p map[string]any) { p["title"] = "  " }, "AD_INVALID_TITLE"},
		{"unknown city", func(p map[string]any) { p["city_id"] = 999999 }, "AD_INVALID_CITY"},
		{"lat without lng", func(p map[string]any) { delete(p, "longitude") }, "AD_INVALID_LOCATION"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload := validPayload(userID)
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
