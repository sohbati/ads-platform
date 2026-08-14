package otptest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	tc "integration-test/internal/testcontainers"
)

const TestMobile = "09125697463"
const OtherMobile = "09111111111"

type ErrorResponse struct {
	Error      string   `json:"error"`
	StatusCode int      `json:"statusCode"`
	Params     []string `json:"params"`
}

type SendOtpResponse struct {
	Message string `json:"message"`
}

type VerifyOtpResponse struct {
	Verified bool   `json:"verified"`
	Message  string `json:"message"`
}

func SetupOtpStack(t *testing.T) (*tc.OtpStack, string, string) {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	t.Cleanup(cancel)

	stack, err := tc.StartOtpStack(ctx, t)
	if err != nil {
		t.Fatalf("start otp stack: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cleanupCancel()
		stack.Terminate(cleanupCtx, t)
	})

	backURL, err := stack.BackURL(ctx)
	if err != nil {
		t.Fatalf("back url: %v", err)
	}

	cacheURL, err := stack.CacheURL(ctx)
	if err != nil {
		t.Fatalf("cache url: %v", err)
	}

	return stack, backURL, cacheURL
}

func SendOTP(ctx context.Context, backURL, mobile string) (int, SendOtpResponse, ErrorResponse) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, backURL+"/api/v1/otp/"+mobile+"/send", nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var success SendOtpResponse
	var failure ErrorResponse

	if resp.StatusCode == http.StatusOK {
		_ = json.Unmarshal(body, &success)
	} else {
		_ = json.Unmarshal(body, &failure)
	}

	return resp.StatusCode, success, failure
}

func VerifyOTP(ctx context.Context, backURL, mobile, otp string) (int, VerifyOtpResponse, ErrorResponse) {
	payload, err := json.Marshal(map[string]string{"otp": otp})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, backURL+"/api/v1/otp/"+mobile+"/verify", bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var success VerifyOtpResponse
	var failure ErrorResponse

	if resp.StatusCode == http.StatusOK {
		_ = json.Unmarshal(body, &success)
	} else {
		_ = json.Unmarshal(body, &failure)
	}

	return resp.StatusCode, success, failure
}

func VerifyOTPRaw(ctx context.Context, backURL, mobile string, body []byte) (int, VerifyOtpResponse, ErrorResponse) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, backURL+"/api/v1/otp/"+mobile+"/verify", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var success VerifyOtpResponse
	var failure ErrorResponse

	if resp.StatusCode == http.StatusOK {
		_ = json.Unmarshal(respBody, &success)
	} else {
		_ = json.Unmarshal(respBody, &failure)
	}

	return resp.StatusCode, success, failure
}

func GetCachedOTP(ctx context.Context, cacheURL, mobile string) (string, error) {
	url := fmt.Sprintf("%s/api/v1/caches/otp/otp:%s", cacheURL, mobile)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cache get status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var otp string
	if err := json.Unmarshal(body, &otp); err != nil {
		return "", err
	}

	return otp, nil
}

func MustSendAndGetOTP(t *testing.T, ctx context.Context, backURL, cacheURL, mobile string) string {
	t.Helper()

	status, resp, errResp := SendOTP(ctx, backURL, mobile)
	if status != http.StatusOK {
		t.Fatalf("send otp: status=%d error=%+v", status, errResp)
	}
	if resp.Message != "otp_sent" {
		t.Fatalf("expected otp_sent, got %q", resp.Message)
	}

	otp, err := GetCachedOTP(ctx, cacheURL, mobile)
	if err != nil {
		t.Fatalf("get cached otp: %v", err)
	}
	if len(otp) != 6 {
		t.Fatalf("expected 6-digit otp, got %q", otp)
	}

	return otp
}
