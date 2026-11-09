// Package smsir implements sms.Sender on the sms.ir REST API (https://sms.ir/rest-api/).
//
// The sandbox uses the same host and endpoints as production. A key created with
// type "Sandbox" in the panel gets simulated responses: no SMS is delivered and no
// credit is charged. In the sandbox only the default verify template is active
// (templateId 123456, parameter "Code").
package smsir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ads-platform-notification/internal/core/sms"
)

const DefaultBaseURL = "https://api.sms.ir"

// Template maps a logical template to an sms.ir verify template.
type Template struct {
	ID int
	// Params maps logical parameter names (sms.ParamCode) to template keys ("Code").
	Params map[string]string
}

type Config struct {
	APIKey  string
	BaseURL string
	// LineNumber is the sender line for free-text (bulk) sends. Zero turns off SendText.
	LineNumber int64
	Templates  map[string]Template
	Timeout    time.Duration
}

type sender struct {
	cfg    Config
	client *http.Client
}

func New(cfg Config) (sms.Sender, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("smsir: API key is required")
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	return &sender{cfg: cfg, client: &http.Client{Timeout: cfg.Timeout}}, nil
}

func (s *sender) Name() string { return "smsir" }

type verifyParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type verifyRequest struct {
	Mobile     string        `json:"mobile"`
	TemplateID int           `json:"templateId"`
	Parameters []verifyParam `json:"parameters"`
}

type bulkRequest struct {
	LineNumber  int64    `json:"lineNumber"`
	MessageText string   `json:"messageText"`
	Mobiles     []string `json:"mobiles"`
}

type envelope struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (s *sender) SendTemplate(ctx context.Context, msg sms.TemplateMessage) (*sms.Result, error) {
	tpl, ok := s.cfg.Templates[msg.Template]
	if !ok {
		return nil, fmt.Errorf("%w: %s/%s", sms.ErrUnknownTemplate, s.Name(), msg.Template)
	}
	req := verifyRequest{Mobile: msg.Mobile, TemplateID: tpl.ID}
	for logical, value := range msg.Params {
		name, ok := tpl.Params[logical]
		if !ok {
			name = logical
		}
		req.Parameters = append(req.Parameters, verifyParam{Name: name, Value: value})
	}

	var data struct {
		MessageID int64 `json:"messageId"`
	}
	if err := s.post(ctx, "/v1/send/verify", req, &data); err != nil {
		return nil, err
	}
	return &sms.Result{Provider: s.Name(), MessageID: strconv.FormatInt(data.MessageID, 10)}, nil
}

func (s *sender) SendText(ctx context.Context, msg sms.TextMessage) (*sms.Result, error) {
	if s.cfg.LineNumber == 0 {
		return nil, fmt.Errorf("%w: %s free text needs SMSIR_LINE_NUMBER", sms.ErrNotSupported, s.Name())
	}
	req := bulkRequest{LineNumber: s.cfg.LineNumber, MessageText: msg.Text, Mobiles: []string{msg.Mobile}}

	var data struct {
		PackID     string   `json:"packId"`
		MessageIDs []*int64 `json:"messageIds"`
	}
	if err := s.post(ctx, "/v1/send/bulk", req, &data); err != nil {
		return nil, err
	}
	if len(data.MessageIDs) == 0 || data.MessageIDs[0] == nil {
		return nil, &sms.ProviderError{Provider: s.Name(), Message: "invalid mobile or text"}
	}
	if *data.MessageIDs[0] == 0 {
		return nil, &sms.ProviderError{Provider: s.Name(), Code: 114, Message: "mobile is blacklisted"}
	}
	return &sms.Result{Provider: s.Name(), MessageID: strconv.FormatInt(*data.MessageIDs[0], 10)}, nil
}

func (s *sender) post(ctx context.Context, path string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.BaseURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-API-KEY", s.cfg.APIKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return &sms.ProviderError{Provider: s.Name(), Message: err.Error(), Retryable: true}
	}
	defer resp.Body.Close()

	var env envelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return &sms.ProviderError{
			Provider: s.Name(), HTTPStatus: resp.StatusCode,
			Message:   fmt.Sprintf("unreadable response (HTTP %d): %v", resp.StatusCode, err),
			Retryable: resp.StatusCode >= 500,
		}
	}
	if resp.StatusCode != http.StatusOK || env.Status != 1 {
		return &sms.ProviderError{
			Provider: s.Name(), HTTPStatus: resp.StatusCode, Code: env.Status,
			Message:   fmt.Sprintf("status %d: %s", env.Status, env.Message),
			Retryable: resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 || env.Status == 0 || env.Status == 20,
		}
	}
	if out != nil && len(env.Data) > 0 && string(env.Data) != "null" {
		if err := json.Unmarshal(env.Data, out); err != nil {
			return fmt.Errorf("smsir: decode data: %w", err)
		}
	}
	return nil
}
