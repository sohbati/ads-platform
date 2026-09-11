package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"ads-platform-ui/internal/business/location/service"
)

const nominatimUserAgent = "ruab.ir/ads-platform-ui"

type NominatimClient struct {
	base        string
	http        *http.Client
	minInterval time.Duration
	mu          sync.Mutex
	lastCall    time.Time
	cityCache   map[string]service.GeoPoint
	revCache    map[string]*service.ReverseResult
}

func NewNominatimClient(baseURL string) *NominatimClient {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		base = "https://nominatim.openstreetmap.org"
	}
	return &NominatimClient{
		base:        base,
		http:        &http.Client{Timeout: 8 * time.Second},
		minInterval: 1100 * time.Millisecond,
		cityCache:   make(map[string]service.GeoPoint),
		revCache:    make(map[string]*service.ReverseResult),
	}
}

func (c *NominatimClient) CityCenter(ctx context.Context, slug, name, lang string) (service.GeoPoint, error) {
	slug = strings.ToLower(strings.TrimSpace(slug))
	if p, ok := service.KnownIranCityCenters[slug]; ok {
		return p, nil
	}
	cacheKey := slug + "|" + strings.TrimSpace(name) + "|" + lang
	c.mu.Lock()
	if cached, ok := c.cityCache[cacheKey]; ok {
		c.mu.Unlock()
		return cached, nil
	}
	c.mu.Unlock()

	q := strings.TrimSpace(name)
	if q == "" {
		q = slug
	}
	if q == "" {
		return service.DefaultCityCenter(), nil
	}

	pt, err := c.search(ctx, q, lang)
	if err != nil || pt.Lat == 0 && pt.Lng == 0 {
		return service.DefaultCityCenter(), nil
	}
	pt.Source = "nominatim"
	c.mu.Lock()
	c.cityCache[cacheKey] = pt
	c.mu.Unlock()
	return pt, nil
}

func (c *NominatimClient) Reverse(ctx context.Context, lat, lng float64, lang string) (*service.ReverseResult, error) {
	key := fmt.Sprintf("%.4f,%.4f,%s", lat, lng, lang)
	c.mu.Lock()
	if cached, ok := c.revCache[key]; ok {
		c.mu.Unlock()
		return cached, nil
	}
	c.mu.Unlock()

	params := url.Values{}
	params.Set("lat", strconv.FormatFloat(lat, 'f', 6, 64))
	params.Set("lon", strconv.FormatFloat(lng, 'f', 6, 64))
	params.Set("format", "json")
	params.Set("addressdetails", "1")
	params.Set("zoom", "16")
	if lang != "" {
		params.Set("accept-language", lang)
	}

	var raw nominatimReverse
	if err := c.getJSON(ctx, "/reverse?"+params.Encode(), &raw); err != nil {
		return nil, err
	}
	out := &service.ReverseResult{
		Lat:          lat,
		Lng:          lng,
		Neighborhood: service.NeighborhoodFromAddress(stringMap(raw.Address)),
		DisplayName:  raw.DisplayName,
	}
	c.mu.Lock()
	c.revCache[key] = out
	c.mu.Unlock()
	return out, nil
}

func (c *NominatimClient) search(ctx context.Context, query, lang string) (service.GeoPoint, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	params.Set("limit", "1")
	params.Set("countrycodes", "ir")
	if lang != "" {
		params.Set("accept-language", lang)
	}
	var hits []nominatimSearchHit
	if err := c.getJSON(ctx, "/search?"+params.Encode(), &hits); err != nil {
		return service.GeoPoint{}, err
	}
	if len(hits) == 0 {
		return service.GeoPoint{}, nil
	}
	lat, err1 := strconv.ParseFloat(hits[0].Lat, 64)
	lng, err2 := strconv.ParseFloat(hits[0].Lon, 64)
	if err1 != nil || err2 != nil {
		return service.GeoPoint{}, nil
	}
	return service.GeoPoint{Lat: lat, Lng: lng}, nil
}

func (c *NominatimClient) getJSON(ctx context.Context, path string, dest any) error {
	c.throttle()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", nominatimUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("nominatim status %d", resp.StatusCode)
	}
	return json.Unmarshal(body, dest)
}

func (c *NominatimClient) throttle() {
	if c.minInterval <= 0 {
		return
	}
	c.mu.Lock()
	wait := c.minInterval - time.Since(c.lastCall)
	if wait > 0 {
		c.mu.Unlock()
		time.Sleep(wait)
		c.mu.Lock()
	}
	c.lastCall = time.Now()
	c.mu.Unlock()
}

type nominatimSearchHit struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type nominatimReverse struct {
	DisplayName string         `json:"display_name"`
	Address     map[string]any `json:"address"`
}

func stringMap(in map[string]any) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		if s, ok := v.(string); ok {
			out[k] = s
		}
	}
	return out
}
