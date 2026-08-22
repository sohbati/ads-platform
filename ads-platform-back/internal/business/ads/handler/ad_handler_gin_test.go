package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/ads/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type fakeAdService struct {
	got service.CreateAdInput
	ad  *model.Ad
	err error
}

func (f *fakeAdService) Create(_ context.Context, in service.CreateAdInput) (*model.Ad, error) {
	f.got = in
	if in.Pictures != nil {
		for _, p := range in.Pictures {
			if p.Body != nil {
				_, _ = io.Copy(io.Discard, p.Body)
			}
		}
	}
	return f.ad, f.err
}

func testRouter(h *AdHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.GlobalErrorHandler())
	r.POST("/api/v1/ads", h.Create)
	return r
}

func TestCreateAdJSONHandler(t *testing.T) {
	fake := &fakeAdService{ad: &model.Ad{ID: 9, Title: "Used laptop", Status: model.AdStatusActive}}
	r := testRouter(NewAdHandler(fake))

	body := []byte(`{
		"user_id":1,"category_id":16,"city_id":1,
		"title":"Used laptop","description":"line1\nline2",
		"latitude":35.6,"longitude":51.4,
		"attrs":{"brand":"lenovo"}
	}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/ads", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if fake.got.UserID != 1 || fake.got.CategoryID != 16 || fake.got.Title != "Used laptop" {
		t.Fatalf("input: %+v", fake.got)
	}
	if fake.got.Description != "line1\nline2" {
		t.Fatalf("description=%q", fake.got.Description)
	}
	if fake.got.Latitude == nil || *fake.got.Latitude != 35.6 {
		t.Fatalf("latitude=%v", fake.got.Latitude)
	}
	if !bytes.Contains(fake.got.Attrs, []byte(`"brand":"lenovo"`)) {
		t.Fatalf("attrs=%s", fake.got.Attrs)
	}
	if len(fake.got.Pictures) != 0 {
		t.Fatalf("expected no pictures, got %d", len(fake.got.Pictures))
	}
}

func TestCreateAdMultipartHandler(t *testing.T) {
	fake := &fakeAdService{ad: &model.Ad{ID: 3, Title: "x"}}
	r := testRouter(NewAdHandler(fake))

	payload, _ := json.Marshal(map[string]any{
		"user_id": 7, "category_id": 16, "city_id": 1,
		"title": "x", "description": "y",
	})
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("payload", string(payload))
	part, err := w.CreateFormFile("pictures", "cover.jpg")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("fake-jpeg"))
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ads", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if fake.got.UserID != 7 {
		t.Fatalf("user_id=%d", fake.got.UserID)
	}
	if len(fake.got.Pictures) != 1 || fake.got.Pictures[0].Filename != "cover.jpg" {
		t.Fatalf("pictures: %+v", fake.got.Pictures)
	}
	if fake.got.Pictures[0].Size != int64(len("fake-jpeg")) {
		t.Fatalf("picture size=%d", fake.got.Pictures[0].Size)
	}
}

func TestCreateAdHandlerServiceError(t *testing.T) {
	fake := &fakeAdService{err: exception.NewAppError("AD_INVALID_TITLE", http.StatusBadRequest)}
	r := testRouter(NewAdHandler(fake))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ads", bytes.NewReader([]byte(`{"user_id":1,"title":"x"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"AD_INVALID_TITLE"`)) {
		t.Fatalf("body=%s", rec.Body.String())
	}
}
