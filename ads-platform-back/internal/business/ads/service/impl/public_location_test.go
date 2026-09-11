package impl

import (
	"math"
	"testing"
)

func TestApproximatePublicCoordsRoundsAndDiffers(t *testing.T) {
	lat, lng := 35.6892, 51.3890
	gotLat, gotLng := approximatePublicCoords(7, lat, lng)
	if gotLat == lat && gotLng == lng {
		t.Fatal("public coords should not equal the exact pin")
	}
	if math.Abs(gotLat-lat) > 0.003 || math.Abs(gotLng-lng) > 0.003 {
		t.Fatalf("moved too far: %v,%v from %v,%v", gotLat, gotLng, lat, lng)
	}
	if gotLat != roundCoord(gotLat) || gotLng != roundCoord(gotLng) {
		t.Fatalf("not rounded: %v %v", gotLat, gotLng)
	}
	aLat, aLng := approximatePublicCoords(7, lat, lng)
	if aLat != gotLat || aLng != gotLng {
		t.Fatal("expected deterministic offset")
	}
}
