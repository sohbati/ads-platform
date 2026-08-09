package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// BackendProxy forwards /api/v1 requests to ads-platform-back.
type BackendProxy struct {
	proxy            *httputil.ReverseProxy
	backendHealthURL string
	client           *http.Client
}

func NewBackendProxy(backendBaseURL string) (*BackendProxy, error) {
	target, err := url.Parse(strings.TrimRight(strings.TrimSpace(backendBaseURL), "/"))
	if err != nil {
		return nil, fmt.Errorf("parse backend url: %w", err)
	}

	rp := httputil.NewSingleHostReverseProxy(target)
	originalDirector := rp.Director
	rp.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
	}
	rp.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, _ error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error":      "BACKEND_UNAVAILABLE",
			"statusCode": http.StatusBadGateway,
		})
	}

	return &BackendProxy{
		proxy:            rp,
		backendHealthURL: target.String() + "/health",
		client:           &http.Client{Timeout: 5 * time.Second},
	}, nil
}

func (p *BackendProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func (p *BackendProxy) BackendHealthy(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.backendHealthURL, nil)
	if err != nil {
		return false
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	return resp.StatusCode == http.StatusOK
}
