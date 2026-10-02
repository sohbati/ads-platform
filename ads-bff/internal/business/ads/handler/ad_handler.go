package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"ads-bff/internal/business/auth/model"
	"ads-bff/internal/business/auth/service"
	backendclient "ads-bff/internal/core/client/backend"
	cacheclient "ads-bff/internal/core/client/cache"
	"ads-bff/internal/core/config"
	"ads-bff/internal/core/exception"
	"ads-bff/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

const maxAdBody = 32 << 20

type AdHandler struct {
	cfg     *config.Config
	auth    service.AuthService
	backend *backendclient.Client
}

func NewAdHandler(cfg *config.Config, auth service.AuthService, backend *backendclient.Client) *AdHandler {
	return &AdHandler{cfg: cfg, auth: auth, backend: backend}
}

func (h *AdHandler) Create(c *gin.Context) {
	customer, err := h.sessionCustomer(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	body, contentType, err := injectCustomerID(c, customer.ID, customer.Mobile)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	status, respBody, err := h.backend.CreateAd(c.Request.Context(), body, contentType)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}

	ct := "application/json"
	c.Data(status, ct, respBody)
}

func (h *AdHandler) ListMine(c *gin.Context) {
	customerID, err := h.sessionCustomerID(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	status, respBody, err := h.backend.GetCustomerAds(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *AdHandler) ListMineStats(c *gin.Context) {
	customerID, err := h.sessionCustomerID(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	status, respBody, err := h.backend.GetCustomerAdStats(c.Request.Context(), customerID, c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *AdHandler) GetMine(c *gin.Context) {
	customerID, err := h.sessionCustomerID(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	adID, err := parsePositiveID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	status, respBody, err := h.backend.GetCustomerAd(c.Request.Context(), customerID, adID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *AdHandler) UpdateMine(c *gin.Context) {
	customer, err := h.sessionCustomer(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	adID, err := parsePositiveID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	body, contentType, err := injectCustomerID(c, customer.ID, customer.Mobile)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	status, respBody, err := h.backend.UpdateCustomerAd(c.Request.Context(), customer.ID, adID, body, contentType)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	c.Data(status, "application/json", respBody)
}

func (h *AdHandler) GetContact(c *gin.Context) {
	if _, err := h.sessionCustomerID(c); err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	adID, err := parsePositiveID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	status, respBody, err := h.backend.GetAdContact(c.Request.Context(), adID)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "BACKEND_UNAVAILABLE", "statusCode": http.StatusBadGateway})
		return
	}
	c.Data(status, "application/json", respBody)
}

func parsePositiveID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest)
	}
	return id, nil
}

func (h *AdHandler) sessionCustomer(c *gin.Context) (*model.SessionCustomer, error) {
	sessionID, err := c.Cookie(h.cfg.SessionCookieName)
	if err != nil || sessionID == "" {
		return nil, exception.NewAppError("AUTH_REQUIRED", http.StatusUnauthorized)
	}

	customer, err := h.auth.GetCurrentCustomer(c.Request.Context(), sessionID)
	if err != nil {
		if cacheclient.IsSessionNotFound(err) {
			return nil, exception.NewAppError("AUTH_REQUIRED", http.StatusUnauthorized)
		}
		return nil, exception.NewAppError("SESSION_LOOKUP_FAILED", http.StatusBadGateway).WithCause(err)
	}
	if customer == nil || customer.ID <= 0 {
		return nil, exception.NewAppError("AUTH_REQUIRED", http.StatusUnauthorized)
	}
	return customer, nil
}

func (h *AdHandler) sessionCustomerID(c *gin.Context) (int64, error) {
	customer, err := h.sessionCustomer(c)
	if err != nil {
		return 0, err
	}
	return customer.ID, nil
}

func injectCustomerID(c *gin.Context, customerID int64, mobile string) ([]byte, string, error) {
	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/form-data") {
		return injectCustomerIDMultipart(c, customerID, mobile)
	}
	return injectCustomerIDJSON(c, customerID, mobile)
}

func injectCustomerIDJSON(c *gin.Context, customerID int64, mobile string) ([]byte, string, error) {
	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, maxAdBody))
	if err != nil {
		return nil, "", exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err)
	}

	payload, err := decodePayload(raw)
	if err != nil {
		return nil, "", err
	}
	payload["customer_id"] = customerID
	ensureContactPhone(payload, mobile)

	out, err := json.Marshal(payload)
	if err != nil {
		return nil, "", exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err)
	}
	return out, "application/json", nil
}

func injectCustomerIDMultipart(c *gin.Context, customerID int64, mobile string) ([]byte, string, error) {
	if err := c.Request.ParseMultipartForm(maxAdBody); err != nil {
		return nil, "", exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err)
	}

	payloadRaw := c.PostForm("payload")
	payload, err := decodePayload([]byte(payloadRaw))
	if err != nil {
		return nil, "", err
	}
	payload["customer_id"] = customerID
	ensureContactPhone(payload, mobile)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, "", exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err)
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if err := w.WriteField("payload", string(payloadJSON)); err != nil {
		return nil, "", exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err)
	}

	form := c.Request.MultipartForm
	if form != nil {
		for _, fh := range form.File["pictures"] {
			if err := copyPicturePart(w, fh); err != nil {
				return nil, "", err
			}
		}
	}
	if err := w.Close(); err != nil {
		return nil, "", exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err)
	}
	return buf.Bytes(), w.FormDataContentType(), nil
}

func copyPicturePart(w *multipart.Writer, fh *multipart.FileHeader) error {
	src, err := fh.Open()
	if err != nil {
		return exception.NewAppError("AD_INVALID_PICTURE", http.StatusBadRequest).WithCause(err)
	}
	defer src.Close()

	header := make(textproto.MIMEHeader)
	filename := filepath.Base(fh.Filename)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="pictures"; filename="%s"`, escapeQuotes(filename)))
	if ct := fh.Header.Get("Content-Type"); ct != "" {
		header.Set("Content-Type", ct)
	}

	part, err := w.CreatePart(header)
	if err != nil {
		return exception.NewAppError("AD_INVALID_PICTURE", http.StatusBadRequest).WithCause(err)
	}
	if _, err := io.Copy(part, src); err != nil {
		return exception.NewAppError("AD_INVALID_PICTURE", http.StatusBadRequest).WithCause(err)
	}
	return nil
}

func decodePayload(raw []byte) (map[string]any, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return nil, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil || payload == nil {
		return nil, exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest)
	}
	delete(payload, "customer_id")
	return payload, nil
}

func ensureContactPhone(payload map[string]any, mobile string) {
	mobile = strings.TrimSpace(mobile)
	if mobile == "" || payload == nil {
		return
	}
	contact, _ := payload["contact"].(map[string]any)
	if contact == nil {
		contact = map[string]any{}
	}
	phone, _ := contact["phone"].(string)
	if strings.TrimSpace(phone) != "" {
		return
	}
	contact["phone"] = mobile
	payload["contact"] = contact
}

func escapeQuotes(s string) string {
	if !utf8.ValidString(s) {
		s = strings.ToValidUTF8(s, "")
	}
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
