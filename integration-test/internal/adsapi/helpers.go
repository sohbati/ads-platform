package adsapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ErrorResponse struct {
	Error      string   `json:"error"`
	StatusCode int      `json:"statusCode"`
	Params     []string `json:"params"`
}

type Ad struct {
	ID          int64           `json:"id"`
	UserID      int64           `json:"user_id"`
	CategoryID  int             `json:"category_id"`
	CityID      int             `json:"city_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	PriceAmount *int64          `json:"price_amount"`
	PriceType   string          `json:"price_type"`
	Currency    string          `json:"currency"`
	Attrs       json.RawMessage `json:"attrs"`
	Media       json.RawMessage `json:"media"`
	Contact     json.RawMessage `json:"contact"`
	Location    json.RawMessage `json:"location"`
}

type PublicAd struct {
	ID           int64           `json:"id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	PriceAmount  *int64          `json:"price_amount"`
	PriceType    string          `json:"price_type"`
	Currency     string          `json:"currency"`
	CityID       int             `json:"city_id"`
	CategoryID   int             `json:"category_id"`
	CityName     string          `json:"city_name,omitempty"`
	Neighborhood string          `json:"neighborhood,omitempty"`
	MapLat       *float64        `json:"map_lat,omitempty"`
	MapLng       *float64        `json:"map_lng,omitempty"`
	Attrs        json.RawMessage `json:"attrs,omitempty"`
	Media        []PublicMedia   `json:"media"`
	PublishedAt  *string         `json:"published_at,omitempty"`
	HasPhone     bool            `json:"has_phone"`
	PhoneMasked  string          `json:"phone_masked,omitempty"`
}

type PublicMedia struct {
	URL     string `json:"url"`
	Thumb   string `json:"thumb"`
	IsCover bool   `json:"is_cover"`
}

type SearchResponse struct {
	Place         string         `json:"place"`
	Category      string         `json:"category"`
	CategoryTitle string         `json:"category_title"`
	Pagination    Pagination     `json:"pagination"`
	Ads           []SearchAdItem `json:"ads"`
}

type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type SearchAdItem struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	CategoryID int    `json:"category_id"`
}

type MediaItem struct {
	ObjectKey   string `json:"object_key"`
	URL         string `json:"url"`
	Thumb       string `json:"thumb"`
	ContentType string `json:"content_type"`
	IsCover     bool   `json:"is_cover"`
}

func PostJSON(ctx context.Context, backURL string, payload any) (int, Ad, ErrorResponse, error) {
	req, err := newJSONRequest(ctx, http.MethodPost, backURL+"/api/v1/ads", payload)
	if err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}
	var ad Ad
	status, failure, err := doJSON(req, http.StatusCreated, &ad)
	return status, ad, failure, err
}

func PostRawJSON(ctx context.Context, backURL string, body []byte) (int, Ad, ErrorResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, backURL+"/api/v1/ads", bytes.NewReader(body))
	if err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	var ad Ad
	status, failure, err := doJSON(req, http.StatusCreated, &ad)
	return status, ad, failure, err
}

func PutJSON(ctx context.Context, backURL string, userID, adID int64, payload any) (int, Ad, ErrorResponse, error) {
	path := backURL + "/api/v1/users/" + strconv.FormatInt(userID, 10) + "/ads/" + strconv.FormatInt(adID, 10)
	req, err := newJSONRequest(ctx, http.MethodPut, path, payload)
	if err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}
	var ad Ad
	status, failure, err := doJSON(req, http.StatusOK, &ad)
	return status, ad, failure, err
}

func GetPublic(ctx context.Context, backURL string, adID int64) (int, PublicAd, ErrorResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, backURL+"/api/v1/ads/"+strconv.FormatInt(adID, 10), nil)
	if err != nil {
		return 0, PublicAd{}, ErrorResponse{}, err
	}
	var ad PublicAd
	status, failure, err := doJSON(req, http.StatusOK, &ad)
	return status, ad, failure, err
}

func Search(ctx context.Context, backURL, place, category, query string) (int, SearchResponse, ErrorResponse, error) {
	u := backURL + "/api/v1/q/" + url.PathEscape(place) + "/" + url.PathEscape(category)
	if query != "" {
		u += "?q=" + url.QueryEscape(query)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, SearchResponse{}, ErrorResponse{}, err
	}
	var out SearchResponse
	status, failure, err := doJSON(req, http.StatusOK, &out)
	return status, out, failure, err
}

func PostMultipart(ctx context.Context, backURL string, payload any, files map[string][]byte) (int, Ad, ErrorResponse, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := w.WriteField("payload", string(raw)); err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}
	for name, data := range files {
		part, err := w.CreateFormFile("pictures", name)
		if err != nil {
			return 0, Ad{}, ErrorResponse{}, err
		}
		if _, err := part.Write(data); err != nil {
			return 0, Ad{}, ErrorResponse{}, err
		}
	}
	if err := w.Close(); err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, backURL+"/api/v1/ads", &buf)
	if err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	var ad Ad
	status, failure, err := doJSON(req, http.StatusCreated, &ad)
	return status, ad, failure, err
}

func newJSONRequest(ctx context.Context, method, url string, payload any) (*http.Request, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func doJSON(req *http.Request, okStatus int, dest any) (int, ErrorResponse, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, ErrorResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, ErrorResponse{}, err
	}

	if resp.StatusCode != okStatus {
		var failure ErrorResponse
		_ = json.Unmarshal(body, &failure)
		return resp.StatusCode, failure, nil
	}

	if dest != nil {
		if err := json.Unmarshal(body, dest); err != nil {
			return resp.StatusCode, ErrorResponse{}, fmt.Errorf("decode: %w; body=%s", err, string(body))
		}
	}
	return resp.StatusCode, ErrorResponse{}, nil
}

func ParseMedia(raw json.RawMessage) ([]MediaItem, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var items []MediaItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func ParseAttrs(raw json.RawMessage) (map[string]any, error) {
	if len(raw) == 0 {
		return map[string]any{}, nil
	}
	var attrs map[string]any
	if err := json.Unmarshal(raw, &attrs); err != nil {
		return nil, err
	}
	if attrs == nil {
		attrs = map[string]any{}
	}
	return attrs, nil
}

// JPEG1x1 is a valid 1×1 JPEG used as a picture upload fixture.
var JPEG1x1 = []byte{
	0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01,
	0x01, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0xff, 0xdb, 0x00, 0x43,
	0x00, 0x08, 0x06, 0x06, 0x07, 0x06, 0x05, 0x08, 0x07, 0x07, 0x07, 0x09,
	0x09, 0x08, 0x0a, 0x0c, 0x14, 0x0d, 0x0c, 0x0b, 0x0b, 0x0c, 0x19, 0x12,
	0x13, 0x0f, 0x14, 0x1d, 0x1a, 0x1f, 0x1e, 0x1d, 0x1a, 0x1c, 0x1c, 0x20,
	0x24, 0x2e, 0x27, 0x20, 0x22, 0x2c, 0x23, 0x1c, 0x1c, 0x28, 0x37, 0x29,
	0x2c, 0x30, 0x31, 0x34, 0x34, 0x34, 0x1f, 0x27, 0x39, 0x3d, 0x38, 0x32,
	0x3c, 0x2e, 0x33, 0x34, 0x32, 0xff, 0xc0, 0x00, 0x0b, 0x08, 0x00, 0x01,
	0x00, 0x01, 0x01, 0x01, 0x11, 0x00, 0xff, 0xc4, 0x00, 0x14, 0x00, 0x01,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x03, 0xff, 0xc4, 0x00, 0x14, 0x10, 0x01, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0xff, 0xda, 0x00, 0x08, 0x01, 0x01, 0x00, 0x00, 0x3f, 0x00,
	0x37, 0xff, 0xd9,
}
