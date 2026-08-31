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
	got     service.CreateAdInput
	ad      *model.Ad
	err     error
	list    []model.UserAdItem
	listErr error
	listID  int64
	gotID   int64
	gotUID  int64
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

func (f *fakeAdService) GetPublic(_ context.Context, adID int64) (*model.PublicAd, error) {
	f.gotID = adID
	if f.err != nil {
		return nil, f.err
	}
	if f.ad == nil {
		return nil, nil
	}
	return &model.PublicAd{ID: f.ad.ID, Title: f.ad.Title, Description: f.ad.Description}, nil
}

func (f *fakeAdService) GetPublicContact(_ context.Context, adID int64) (*model.PublicContact, error) {
	f.gotID = adID
	if f.err != nil {
		return nil, f.err
	}
	return &model.PublicContact{Phone: "09121110001"}, nil
}

func (f *fakeAdService) GetForOwner(_ context.Context, userID, adID int64) (*model.Ad, error) {
	f.gotUID = userID
	f.gotID = adID
	return f.ad, f.err
}

func (f *fakeAdService) Update(_ context.Context, adID int64, in service.CreateAdInput) (*model.Ad, error) {
	f.gotID = adID
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

func (f *fakeAdService) ListByUser(_ context.Context, userID int64) ([]model.UserAdItem, error) {
	f.listID = userID
	return f.list, f.listErr
}

func testRouter(h *AdHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.GlobalErrorHandler())
	r.POST("/api/v1/ads", h.Create)
	r.GET("/api/v1/ads/:id", h.GetPublic)
	r.GET("/api/v1/ads/:id/contact", h.GetPublicContact)
	r.GET("/api/v1/users/:userId/ads", h.ListByUser)
	r.GET("/api/v1/users/:userId/ads/:adId", h.GetForOwner)
	r.PUT("/api/v1/users/:userId/ads/:adId", h.Update)
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

func TestListByUserHandler(t *testing.T) {
	fake := &fakeAdService{list: []model.UserAdItem{{ID: 3, Title: "Mine"}}}
	r := testRouter(NewAdHandler(fake))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/42/ads", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if fake.listID != 42 {
		t.Fatalf("userID=%d", fake.listID)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"Mine"`)) {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestListByUserHandlerInvalidID(t *testing.T) {
	r := testRouter(NewAdHandler(&fakeAdService{}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/abc/ads", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestGetPublicHandler(t *testing.T) {
	fake := &fakeAdService{ad: &model.Ad{ID: 9, Title: "Used laptop"}}
	r := testRouter(NewAdHandler(fake))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ads/9", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if fake.gotID != 9 {
		t.Fatalf("id=%d", fake.gotID)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"title":"Used laptop"`)) {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestGetPublicContactHandler(t *testing.T) {
	fake := &fakeAdService{ad: &model.Ad{ID: 9}}
	r := testRouter(NewAdHandler(fake))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ads/9/contact", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if fake.gotID != 9 {
		t.Fatalf("id=%d", fake.gotID)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"phone":"09121110001"`)) {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestGetForOwnerHandler(t *testing.T) {
	fake := &fakeAdService{ad: &model.Ad{ID: 9, Title: "Mine"}}
	r := testRouter(NewAdHandler(fake))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/42/ads/9", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if fake.gotUID != 42 || fake.gotID != 9 {
		t.Fatalf("ids user=%d ad=%d", fake.gotUID, fake.gotID)
	}
}

func TestUpdateHandlerUsesPathUser(t *testing.T) {
	fake := &fakeAdService{ad: &model.Ad{ID: 9, Title: "Updated"}}
	r := testRouter(NewAdHandler(fake))

	body := []byte(`{"user_id":99,"category_id":13,"city_id":1,"title":"Updated","description":"ok"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/42/ads/9", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if fake.got.UserID != 42 || fake.gotID != 9 || fake.got.Title != "Updated" {
		t.Fatalf("input: %+v id=%d", fake.got, fake.gotID)
	}
}
