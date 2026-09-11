package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ads-platform-ui/internal/business/location/service"

	"github.com/gin-gonic/gin"
)

type fakeGeo struct {
	point service.GeoPoint
	rev   *service.ReverseResult
	err   error
}

func (f fakeGeo) CityCenter(context.Context, string, string, string) (service.GeoPoint, error) {
	return f.point, f.err
}
func (f fakeGeo) Reverse(context.Context, float64, float64, string) (*service.ReverseResult, error) {
	return f.rev, f.err
}

func TestGeoReverseRejectsBadCoords(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewGeoHandler(fakeGeo{})
	r.GET("/reverse", h.Reverse)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reverse?lat=200&lng=51", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status %d", w.Code)
	}
}

func TestGeoReverseOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewGeoHandler(fakeGeo{rev: &service.ReverseResult{Lat: 35.7, Lng: 51.4, Neighborhood: "Vanak"}})
	r.GET("/reverse", h.Reverse)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reverse?lat=35.7&lng=51.4&lang=fa", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	var got service.ReverseResult
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil || got.Neighborhood != "Vanak" {
		t.Fatalf("%v %+v", err, got)
	}
}
