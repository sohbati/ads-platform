package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		http:    httpClient,
	}
}

func NewAdsClient(baseURL string) *Client {
	return NewClient(baseURL, &http.Client{Timeout: 60 * time.Second})
}

type VerifyOtpRequest struct {
	Otp string `json:"otp"`
}

type VerifyOtpResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"statusCode"`
}

type User struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	NationalId string `json:"national_id"`
}

func (c *Client) SendOTP(ctx context.Context, mobile string) (int, []byte, error) {
	url := fmt.Sprintf("%s/api/v1/otp/%s/send", c.baseURL, mobile)
	return c.post(ctx, url, nil)
}

func (c *Client) VerifyOTP(ctx context.Context, mobile, otp string) (*VerifyOtpResponse, int, []byte, error) {
	body, err := json.Marshal(VerifyOtpRequest{Otp: otp})
	if err != nil {
		return nil, 0, nil, err
	}

	url := fmt.Sprintf("%s/api/v1/otp/%s/verify", c.baseURL, mobile)
	status, respBody, err := c.post(ctx, url, body)
	if err != nil {
		return nil, status, respBody, err
	}
	if status != http.StatusOK {
		return nil, status, respBody, nil
	}

	var resp VerifyOtpResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, status, respBody, fmt.Errorf("parse verify response: %w", err)
	}
	return &resp, status, respBody, nil
}

func (c *Client) GetUserByMobile(ctx context.Context, mobile string) (*User, int, []byte, error) {
	url := fmt.Sprintf("%s/api/v1/users/mobile/%s", c.baseURL, mobile)
	status, respBody, err := c.get(ctx, url)
	if err != nil {
		return nil, status, respBody, err
	}
	if status != http.StatusOK {
		return nil, status, respBody, nil
	}

	var user User
	if err := json.Unmarshal(respBody, &user); err != nil {
		return nil, status, respBody, fmt.Errorf("parse user response: %w", err)
	}
	return &user, status, respBody, nil
}

func (c *Client) RegisterUserByMobile(ctx context.Context, mobile string) (*User, int, []byte, error) {
	url := fmt.Sprintf("%s/api/v1/users/register-by-mobile/%s", c.baseURL, mobile)
	status, respBody, err := c.post(ctx, url, nil)
	if err != nil {
		return nil, status, respBody, err
	}
	if status != http.StatusOK {
		return nil, status, respBody, nil
	}

	var user User
	if err := json.Unmarshal(respBody, &user); err != nil {
		return nil, status, respBody, fmt.Errorf("parse register response: %w", err)
	}
	return &user, status, respBody, nil
}

func (c *Client) GetAdContact(ctx context.Context, adID int64) (int, []byte, error) {
	url := fmt.Sprintf("%s/api/v1/ads/%d/contact", c.baseURL, adID)
	return c.get(ctx, url)
}

func (c *Client) GetUserProfile(ctx context.Context, userID int64) (int, []byte, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d/profile", c.baseURL, userID)
	return c.get(ctx, url)
}

func (c *Client) GetUserAds(ctx context.Context, userID int64) (int, []byte, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d/ads", c.baseURL, userID)
	return c.get(ctx, url)
}

func (c *Client) GetUserAdStats(ctx context.Context, userID int64, from, to string) (int, []byte, error) {
	u, err := url.Parse(fmt.Sprintf("%s/api/v1/users/%d/ad-stats", c.baseURL, userID))
	if err != nil {
		return 0, nil, err
	}
	q := u.Query()
	if from != "" {
		q.Set("from", from)
	}
	if to != "" {
		q.Set("to", to)
	}
	u.RawQuery = q.Encode()
	return c.get(ctx, u.String())
}

func (c *Client) GetUserAd(ctx context.Context, userID, adID int64) (int, []byte, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d/ads/%d", c.baseURL, userID, adID)
	return c.get(ctx, url)
}

func (c *Client) UpdateUserAd(ctx context.Context, userID, adID int64, body []byte, contentType string) (int, []byte, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d/ads/%d", c.baseURL, userID, adID)
	reader := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, reader)
	if err != nil {
		return 0, nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}

func (c *Client) PutUserProfile(ctx context.Context, userID int64, body []byte) (int, []byte, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d/profile", c.baseURL, userID)
	return c.put(ctx, url, body)
}

func (c *Client) put(ctx context.Context, url string, body []byte) (int, []byte, error) {
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, reader)
	if err != nil {
		return 0, nil, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}

func (c *Client) CreateAd(ctx context.Context, body []byte, contentType string) (int, []byte, error) {
	url := c.baseURL + "/api/v1/ads"
	reader := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reader)
	if err != nil {
		return 0, nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}

func (c *Client) post(ctx context.Context, url string, body []byte) (int, []byte, error) {
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reader)
	if err != nil {
		return 0, nil, err
	}
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}

func (c *Client) get(ctx context.Context, url string) (int, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, data, nil
}
