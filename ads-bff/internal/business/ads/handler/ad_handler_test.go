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

	"ads-bff/internal/business/auth/model"
	cacheclient "ads-bff/internal/core/client/backend"
	cacheerr "ads-bff/internal/core/client/cache"
	"ads-bff/internal/core/config"
	"ads-bff/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type fakeAuth struct {
	user *model.SessionUser
	err  error
}

func (f *fakeAuth) SendOTP(context.Context, string) (int, []byte, error) {
	return 0, nil, nil
}
func (f *fakeAuth) LoginWithOTP(context.Context, string, string) (*model.LoginResponse, string, int, []byte, error) {
	return nil, "", 0, nil, nil
}
func (f *fakeAuth) GetCurrentUser(context.Context, string) (*model.SessionUser, error) {
	return f.user, f.err
}
func (f *fakeAuth) Logout(context.Context, string) error { return nil }

func testAdsRouter(h *AdHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.GlobalErrorHandler())
	r.POST("/api/v1/ads", h.Create)
	r.GET("/api/v1/ads/:id/contact", h.GetContact)
	r.GET("/api/v1/me/ads", h.ListMine)
	r.GET("/api/v1/me/ads/:id", h.GetMine)
	r.PUT("/api/v1/me/ads/:id", h.UpdateMine)
	return r
}

func newHandler(t *testing.T, auth *fakeAuth, backend http.HandlerFunc) *AdHandler {
	t.Helper()
	srv := httptest.NewServer(backend)
	t.Cleanup(srv.Close)
	cfg := &config.Config{SessionCookieName: "ads_session"}
	return NewAdHandler(cfg, auth, cacheclient.NewAdsClient(srv.URL))
}

func withSession(req *http.Request) {
	req.AddCookie(&http.Cookie{Name: "ads_session", Value: "sess-1"})
}

func TestCreateAdRequiresAuth(t *testing.T) {
	h := newHandler(t, &fakeAuth{err: cacheerr.ErrSessionNotFound}, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("backend should not be called")
	})
	r := testAdsRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ads", bytes.NewReader([]byte(`{"title":"x"}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("AUTH_REQUIRED")) {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestCreateAdJSONInjectsSessionUserID(t *testing.T) {
	var got map[string]any
	h := newHandler(t, &fakeAuth{user: &model.SessionUser{ID: 42}}, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &got); err != nil {
			t.Errorf("backend body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":9,"user_id":42}`))
	})
	r := testAdsRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ads", bytes.NewReader([]byte(`{
		"user_id":99,"category_id":16,"city_id":1,"title":"Used laptop","description":"ok"
	}`)))
	req.Header.Set("Content-Type", "application/json")
	withSession(req)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if got["user_id"] != float64(42) {
		t.Fatalf("backend user_id=%v, want 42 (not client 99)", got["user_id"])
	}
	if got["title"] != "Used laptop" || got["category_id"] != float64(16) {
		t.Fatalf("payload: %+v", got)
	}
}

func TestCreateAdMultipartInjectsSessionUserID(t *testing.T) {
	var got map[string]any
	var pictures int
	h := newHandler(t, &fakeAuth{user: &model.SessionUser{ID: 7}}, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Errorf("parse: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal([]byte(r.FormValue("payload")), &got); err != nil {
			t.Errorf("payload: %v", err)
		}
		pictures = len(r.MultipartForm.File["pictures"])
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":3}`))
	})
	r := testAdsRouter(h)

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("payload", `{"user_id":1,"category_id":16,"city_id":1,"title":"pic","description":"d"}`)
	part, err := mw.CreateFormFile("pictures", "cover.jpg")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("jpeg-bytes"))
	_ = mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ads", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	withSession(req)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if got["user_id"] != float64(7) {
		t.Fatalf("backend user_id=%v", got["user_id"])
	}
	if pictures != 1 {
		t.Fatalf("pictures=%d", pictures)
	}
}

func TestListMineRequiresAuth(t *testing.T) {
	h := newHandler(t, &fakeAuth{err: cacheerr.ErrSessionNotFound}, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("backend should not be called")
	})
	r := testAdsRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me/ads", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestListMineUsesSessionUserID(t *testing.T) {
	var gotPath string
	h := newHandler(t, &fakeAuth{user: &model.SessionUser{ID: 42}}, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ads":[{"id":9,"title":"Mine"}]}`))
	})
	r := testAdsRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me/ads", nil)
	withSession(req)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if gotPath != "/api/v1/users/42/ads" {
		t.Fatalf("backend path=%s", gotPath)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"Mine"`)) {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestGetMineUsesSessionUserID(t *testing.T) {
	var gotPath string
	h := newHandler(t, &fakeAuth{user: &model.SessionUser{ID: 42}}, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"id":9,"title":"Mine"}`))
	})
	r := testAdsRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/me/ads/9", nil)
	withSession(req)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if gotPath != "/api/v1/users/42/ads/9" {
		t.Fatalf("backend path=%s", gotPath)
	}
}

func TestUpdateMineUsesSessionUserID(t *testing.T) {
	var gotPath string
	var got map[string]any
	h := newHandler(t, &fakeAuth{user: &model.SessionUser{ID: 42}}, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if r.Method != http.MethodPut {
			t.Errorf("method=%s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		_, _ = w.Write([]byte(`{"id":9,"title":"Edited"}`))
	})
	r := testAdsRouter(h)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/me/ads/9", bytes.NewReader([]byte(`{"user_id":99,"title":"Edited","description":"ok","category_id":13,"city_id":1}`)))
	req.Header.Set("Content-Type", "application/json")
	withSession(req)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if gotPath != "/api/v1/users/42/ads/9" {
		t.Fatalf("backend path=%s", gotPath)
	}
	if got["user_id"] != float64(42) {
		t.Fatalf("backend user_id=%v", got["user_id"])
	}
}

func TestGetContactRequiresAuth(t *testing.T) {
	h := newHandler(t, &fakeAuth{err: cacheerr.ErrSessionNotFound}, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("backend should not be called")
	})
	r := testAdsRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ads/9/contact", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("AUTH_REQUIRED")) {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestGetContactProxiesWhenAuthenticated(t *testing.T) {
	var gotPath string
	h := newHandler(t, &fakeAuth{user: &model.SessionUser{ID: 42}}, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"phone":"09121110001"}`))
	})
	r := testAdsRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ads/9/contact", nil)
	withSession(req)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if gotPath != "/api/v1/ads/9/contact" {
		t.Fatalf("backend path=%s", gotPath)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(`"phone":"09121110001"`)) {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestCreateAdJSONInjectsSessionMobileWhenContactPhoneMissing(t *testing.T) {
	var got map[string]any
	h := newHandler(t, &fakeAuth{user: &model.SessionUser{ID: 42, Mobile: "+989121110001"}}, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &got); err != nil {
			t.Errorf("backend body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":9,"user_id":42}`))
	})
	r := testAdsRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/ads", bytes.NewReader([]byte(`{
		"category_id":16,"city_id":1,"title":"Used laptop","description":"ok","contact":{"chat_enabled":true}
	}`)))
	req.Header.Set("Content-Type", "application/json")
	withSession(req)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	contact, _ := got["contact"].(map[string]any)
	if contact == nil || contact["phone"] != "+989121110001" {
		t.Fatalf("contact=%v", got["contact"])
	}
}
