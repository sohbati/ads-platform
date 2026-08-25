package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cache-service/internal/simplecache/attrschema/model"
)

type CDNAttrSchemaClient interface {
	ListAttrSchemas(ctx context.Context) ([]model.AttrSchema, error)
}

type httpCDNAttrSchemaClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewCDNAttrSchemaClient(baseURL string) CDNAttrSchemaClient {
	return &httpCDNAttrSchemaClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *httpCDNAttrSchemaClient) ListAttrSchemas(ctx context.Context) ([]model.AttrSchema, error) {
	url := c.baseURL + "/api/attr-schemas"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cdn attr-schemas status %d: %s", resp.StatusCode, string(body))
	}

	var items []model.AttrSchema
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	return items, nil
}
