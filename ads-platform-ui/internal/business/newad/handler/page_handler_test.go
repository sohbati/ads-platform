package handler

import (
	"encoding/json"
	"testing"
)

func TestCoordsFromLocation(t *testing.T) {
	raw := json.RawMessage(`{"lat":35.7,"lng":51.4,"neighborhood":"Vanak"}`)
	if neighborhoodFromLocation(raw) != "Vanak" {
		t.Fatal(neighborhoodFromLocation(raw))
	}
	lat := latFromLocation(raw)
	lng := lngFromLocation(raw)
	if lat == nil || *lat != 35.7 || lng == nil || *lng != 51.4 {
		t.Fatalf("lat=%v lng=%v", lat, lng)
	}
	if latFromLocation(json.RawMessage(`{}`)) != nil {
		t.Fatal("empty")
	}
}
