package service

// KnownIranCityCenters are fallback map centers keyed by ads-platform-cdn city slug.
// Values are approximate city centroids, not listing pins.
var KnownIranCityCenters = map[string]GeoPoint{
	"tehran":       {Lat: 35.6892, Lng: 51.3890, Source: "catalog"},
	"mashhad":      {Lat: 36.2970, Lng: 59.6062, Source: "catalog"},
	"karaj":        {Lat: 35.8327, Lng: 50.9915, Source: "catalog"},
	"shiraz":       {Lat: 29.5918, Lng: 52.5836, Source: "catalog"},
	"isfahan":      {Lat: 32.6546, Lng: 51.6680, Source: "catalog"},
	"ahvaz":        {Lat: 31.3183, Lng: 48.6706, Source: "catalog"},
	"tabriz":       {Lat: 38.0800, Lng: 46.2919, Source: "catalog"},
	"kermanshah":   {Lat: 34.3142, Lng: 47.0650, Source: "catalog"},
	"qom":          {Lat: 34.6401, Lng: 50.8764, Source: "catalog"},
	"rasht":        {Lat: 37.2808, Lng: 49.5832, Source: "catalog"},
	"kerman":       {Lat: 30.2839, Lng: 57.0834, Source: "catalog"},
	"urmia":        {Lat: 37.5527, Lng: 45.0761, Source: "catalog"},
	"zahedan":      {Lat: 29.4963, Lng: 60.8629, Source: "catalog"},
	"hamadan":      {Lat: 34.7992, Lng: 48.5146, Source: "catalog"},
	"yazd":         {Lat: 31.8974, Lng: 54.3676, Source: "catalog"},
	"ardabil":      {Lat: 38.2498, Lng: 48.2933, Source: "catalog"},
	"bandar-abbas": {Lat: 27.1832, Lng: 56.2666, Source: "catalog"},
	"arak":         {Lat: 34.0917, Lng: 49.6892, Source: "catalog"},
	"eslamshahr":   {Lat: 35.5446, Lng: 51.2300, Source: "catalog"},
	"zanjan":       {Lat: 36.6736, Lng: 48.4787, Source: "catalog"},
}

func DefaultCityCenter() GeoPoint {
	return KnownIranCityCenters["tehran"]
}
