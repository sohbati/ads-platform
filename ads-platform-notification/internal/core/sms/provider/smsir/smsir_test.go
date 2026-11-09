package smsir

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"ads-platform-notification/internal/core/sms"
)

func newTestSender(t *testing.T, h http.HandlerFunc) sms.Sender {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	s, err := New(Config{
		APIKey:  "test-key",
		BaseURL: srv.URL,
		Templates: map[string]Template{
			sms.TemplateOTP: {ID: 123456, Params: map[string]string{sms.ParamCode: "Code"}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSendTemplateMapsToVerify(t *testing.T) {
	s := newTestSender(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/send/verify" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-API-KEY") != "test-key" {
			t.Errorf("missing api key header")
		}
		var req verifyRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Mobile != "09120000000" || req.TemplateID != 123456 ||
			len(req.Parameters) != 1 || req.Parameters[0] != (verifyParam{"Code", "4821"}) {
			t.Errorf("unexpected body %+v", req)
		}
		w.Write([]byte(`{"status":1,"message":"ok","data":{"messageId":89545112,"cost":1.0}}`))
	})

	res, err := s.SendTemplate(context.Background(), sms.TemplateMessage{
		Mobile: "09120000000", Template: sms.TemplateOTP, Params: map[string]string{sms.ParamCode: "4821"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.MessageID != "89545112" || res.Provider != "smsir" {
		t.Errorf("unexpected result %+v", res)
	}
}

func TestProviderErrorIsReported(t *testing.T) {
	s := newTestSender(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"data":null,"status":10,"message":"invalid key"}`))
	})

	_, err := s.SendTemplate(context.Background(), sms.TemplateMessage{
		Mobile: "09120000000", Template: sms.TemplateOTP, Params: map[string]string{sms.ParamCode: "1"},
	})
	var pe *sms.ProviderError
	if !errors.As(err, &pe) || pe.Code != 10 || pe.HTTPStatus != 401 || pe.Retryable {
		t.Fatalf("unexpected error %#v", err)
	}
}

func TestUnknownTemplate(t *testing.T) {
	s := newTestSender(t, func(w http.ResponseWriter, r *http.Request) { t.Error("should not call API") })
	_, err := s.SendTemplate(context.Background(), sms.TemplateMessage{Mobile: "0912", Template: "nope"})
	if !errors.Is(err, sms.ErrUnknownTemplate) {
		t.Fatalf("got %v", err)
	}
}
