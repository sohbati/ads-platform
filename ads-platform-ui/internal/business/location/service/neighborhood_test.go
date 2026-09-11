package service

import "testing"

func TestNeighborhoodFromAddress(t *testing.T) {
	got := NeighborhoodFromAddress(map[string]string{
		"city":           "Tehran",
		"suburb":         "Vanak",
		"neighbourhood":  "Vanak Sq",
		"country":        "Iran",
	})
	if got != "Vanak Sq" {
		t.Fatalf("got %q", got)
	}
	if NeighborhoodFromAddress(nil) != "" {
		t.Fatal("nil")
	}
	if NeighborhoodFromAddress(map[string]string{"city": "Tehran"}) != "" {
		t.Fatal("city is not a neighborhood")
	}
}
