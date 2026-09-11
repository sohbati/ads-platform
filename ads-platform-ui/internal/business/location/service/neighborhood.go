package service

import "strings"

// NeighborhoodFromAddress picks the most specific OSM address field for a listing.
func NeighborhoodFromAddress(addr map[string]string) string {
	if addr == nil {
		return ""
	}
	keys := []string{
		"neighbourhood",
		"neighborhood",
		"suburb",
		"quarter",
		"city_district",
		"district",
		"village",
		"town",
		"hamlet",
	}
	for _, key := range keys {
		if v := strings.TrimSpace(addr[key]); v != "" {
			return v
		}
	}
	return ""
}
