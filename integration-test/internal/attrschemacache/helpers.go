package attrschemacache

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AttrSchema struct {
	Name       string          `json:"name"`
	Title      string          `json:"title"`
	JSONSchema json.RawMessage `json:"jsonSchema"`
}

type ErrorResponse struct {
	Error      string   `json:"error"`
	StatusCode int      `json:"statusCode"`
	Params     []string `json:"params"`
}

func GetAttrSchemasByNames(ctx context.Context, cacheURL string, names ...string) (int, []AttrSchema, ErrorResponse, error) {
	u := cacheURL + "/api/v1/caches/attr-schemas/by-names?names=" + url.QueryEscape(strings.Join(names, ","))
	return getAttrSchemas(ctx, u)
}

func getAttrSchemas(ctx context.Context, u string) (int, []AttrSchema, ErrorResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, nil, ErrorResponse{}, err
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, ErrorResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, ErrorResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		var failure ErrorResponse
		_ = json.Unmarshal(body, &failure)
		return resp.StatusCode, nil, failure, nil
	}

	var items []AttrSchema
	if err := json.Unmarshal(body, &items); err != nil {
		return resp.StatusCode, nil, ErrorResponse{}, fmt.Errorf("decode attr schemas: %w; body=%s", err, string(body))
	}
	return resp.StatusCode, items, ErrorResponse{}, nil
}
