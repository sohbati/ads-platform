package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ads-bff/internal/business/auth/model"
	backendclient "ads-bff/internal/core/client/backend"
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

func testRouter(h *ProfileHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.GlobalErrorHandler())
	r.GET("/api/v1/me/profile", h.Get)
	r.PUT("/api/v1/me/profile", h.Put)
	return r
}

func newHandler(t *testing.T, auth *fakeAuth, backend http.HandlerFunc) *ProfileHandler {
	t.Helper()
	srv := httptest.NewServer(backend)
	t.Cleanup(srv.Close)
	cfg := &config.Config{SessionCookieName: "ads_session"}
	return NewProfileHandler(cfg, auth, backendclient.NewClient(srv.URL, nil))
}

func withSession(req *http.Request) {
	req.AddCookie(&http.Cookie{Name: "ads_session", Value: "sess-1"})
}

func TestGetProfileRequiresAuth(t *testing.T) {
	h := newHandler(t, &fakeAuth{err: cacheerr.ErrSessionNotFound}, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("backend should not be called")
	})
	r := testRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me/profile", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestGetProfileUsesSessionUserID(t *testing.T) {
	var path string
	h := newHandler(t, &fakeAuth{user: &model.SessionUser{ID: 42}}, func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user_id":42,"location_slugs":["tehran"]}`))
	})
	r := testRouter(h)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/me/profile", nil)
	withSession(req)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if path != "/api/v1/users/42/profile" {
		t.Fatalf("path=%s", path)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte("tehran")) {
		t.Fatalf("body=%s", rec.Body.String())
	}
}

func TestPutProfileForwardsBodyToSessionUser(t *testing.T) {
	var path, method string
	var body []byte
	h := newHandler(t, &fakeAuth{user: &model.SessionUser{ID: 7}}, func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		body, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"user_id":7,"location_slugs":["karaj"]}`))
	})
	r := testRouter(h)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/me/profile", strings.NewReader(`{"location_slugs":["karaj"]}`))
	req.Header.Set("Content-Type", "application/json")
	withSession(req)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if method != http.MethodPut || path != "/api/v1/users/7/profile" {
		t.Fatalf("method=%s path=%s", method, path)
	}
	if !bytes.Contains(body, []byte("karaj")) {
		t.Fatalf("backend body=%s", body)
	}
}
