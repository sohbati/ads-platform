package adsapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

type MediaItem struct {
	ObjectKey   string `json:"object_key"`
	URL         string `json:"url"`
	Thumb       string `json:"thumb"`
	ContentType string `json:"content_type"`
	IsCover     bool   `json:"is_cover"`
}

func PostJSON(ctx context.Context, backURL string, payload any) (int, Ad, ErrorResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, backURL+"/api/v1/ads", bytes.NewReader(body))
	if err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	return doCreate(req)
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
	return doCreate(req)
}

func doCreate(req *http.Request) (int, Ad, ErrorResponse, error) {
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, Ad{}, ErrorResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, Ad{}, ErrorResponse{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		var failure ErrorResponse
		_ = json.Unmarshal(body, &failure)
		return resp.StatusCode, Ad{}, failure, nil
	}

	var ad Ad
	if err := json.Unmarshal(body, &ad); err != nil {
		return resp.StatusCode, Ad{}, ErrorResponse{}, fmt.Errorf("decode ad: %w; body=%s", err, string(body))
	}
	return resp.StatusCode, ad, ErrorResponse{}, nil
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
