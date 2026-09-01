package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ads-bff/internal/business/auth/model"
	cacheerr "ads-bff/internal/core/client/cache"
	"ads-bff/internal/core/config"
	"ads-bff/internal/core/middleware"
	"ads-bff/internal/core/ratelimit"

	"github.com/gin-gonic/gin"
)

type fakeAuth struct {
	user *model.SessionUser
	err  error
}

func (f *fakeAuth) SendOTP(context.Context, string) (int, []byte, error) { return 0, nil, nil }
func (f *fakeAuth) LoginWithOTP(context.Context, string, string) (*model.LoginResponse, string, int, []byte, error) {
	return nil, "", 0, nil, nil
}
func (f *fakeAuth) GetCurrentUser(context.Context, string) (*model.SessionUser, error) {
	return f.user, f.err
}
func (f *fakeAuth) Logout(context.Context, string) error { return nil }

type fakePub struct {
	payloads [][]byte
	err      error
}

func (f *fakePub) Publish(_ context.Context, payload []byte) error {
	f.payloads = append(f.payloads, append([]byte(nil), payload...))
	return f.err
}

func testStatsRouter(h *StatsHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.GlobalErrorHandler())
	r.POST("/api/v1/stats/events", h.Ingest)
	return r
}

func validBody() []byte {
	return []byte(`{"ad_id":55,"event":"view","viewer_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","occurred_at":"2026-08-30T10:00:00Z"}`)
}

func TestIngestRejectsInvalidBody(t *testing.T) {
	h := NewStatsHandler(&config.Config{SessionCookieName: "ads_session"}, &fakeAuth{}, &fakePub{}, ratelimit.New(0, 1000))
	r := testStatsRouter(h)

	cases := []string{
		`{}`,
		`{"ad_id":0,"event":"view","viewer_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","occurred_at":"2026-08-30T10:00:00Z"}`,
		`{"ad_id":55,"event":"click","viewer_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","occurred_at":"2026-08-30T10:00:00Z"}`,
		`{"ad_id":55,"event":"view","viewer_id":"not-a-uuid","occurred_at":"2026-08-30T10:00:00Z"}`,
		`{"ad_id":55,"event":"view","viewer_id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","occurred_at":"yesterday"}`,
	}
	for _, body := range cases {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/stats/events", bytes.NewReader([]byte(body)))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("body=%s status=%d", body, rec.Code)
		}
	}
}

func TestIngestReturns204AndPublishes(t *testing.T) {
	pub := &fakePub{}
	h := NewStatsHandler(&config.Config{SessionCookieName: "ads_session"}, &fakeAuth{}, pub, ratelimit.New(0, 1000))
	r := testStatsRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/stats/events", bytes.NewReader(validBody()))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if len(pub.payloads) != 1 {
		t.Fatalf("published=%d", len(pub.payloads))
	}
	var got map[string]any
	if err := json.Unmarshal(pub.payloads[0], &got); err != nil {
		t.Fatal(err)
	}
	if got["ad_id"] != float64(55) || got["event"] != "view" {
		t.Fatalf("payload=%s", pub.payloads[0])
	}
	if _, ok := got["session_user_id"]; ok {
		t.Fatalf("guest should omit session_user_id: %s", pub.payloads[0])
	}
}

func TestIngestAttachesSessionUserID(t *testing.T) {
	pub := &fakePub{}
	h := NewStatsHandler(&config.Config{SessionCookieName: "ads_session"}, &fakeAuth{user: &model.SessionUser{ID: 42}}, pub, ratelimit.New(0, 1000))
	r := testStatsRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/stats/events", bytes.NewReader(validBody()))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "ads_session", Value: "sess-1"})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status=%d", rec.Code)
	}
	if !bytes.Contains(pub.payloads[0], []byte(`"session_user_id":42`)) {
		t.Fatalf("payload=%s", pub.payloads[0])
	}
}

func TestIngestReturns204WhenPublishFails(t *testing.T) {
	pub := &fakePub{err: context.DeadlineExceeded}
	h := NewStatsHandler(&config.Config{SessionCookieName: "ads_session"}, &fakeAuth{err: cacheerr.ErrSessionNotFound}, pub, ratelimit.New(0, 1000))
	r := testStatsRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/stats/events", bytes.NewReader(validBody()))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
