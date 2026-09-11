package impl

import "math"

// publicCoordPlaces rounds published map points to ~100 m so a listing does not
// reveal the exact pin the seller dropped.
const publicCoordPlaces = 3

func approximatePublicCoords(adID int64, lat, lng float64) (float64, float64) {
	h := uint64(adID) * 2654435761
	dlat := (float64(int(h%31)) - 15) / 10000
	dlng := (float64(int((h/31)%31)) - 15) / 10000
	return roundCoord(lat + dlat), roundCoord(lng + dlng)
}

func roundCoord(v float64) float64 {
	p := math.Pow(10, publicCoordPlaces)
	return math.Round(v*p) / p
}
