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
	"strings"
	"unicode/utf8"

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
	userID, err := h.sessionUserID(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	body, contentType, err := injectUserID(c, userID)
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

func (h *AdHandler) sessionUserID(c *gin.Context) (int64, error) {
	sessionID, err := c.Cookie(h.cfg.SessionCookieName)
	if err != nil || sessionID == "" {
		return 0, exception.NewAppError("AUTH_REQUIRED", http.StatusUnauthorized)
	}

	user, err := h.auth.GetCurrentUser(c.Request.Context(), sessionID)
	if err != nil {
		if cacheclient.IsSessionNotFound(err) {
			return 0, exception.NewAppError("AUTH_REQUIRED", http.StatusUnauthorized)
		}
		return 0, exception.NewAppError("SESSION_LOOKUP_FAILED", http.StatusBadGateway).WithCause(err)
	}
	if user == nil || user.ID <= 0 {
		return 0, exception.NewAppError("AUTH_REQUIRED", http.StatusUnauthorized)
	}
	return user.ID, nil
}

func injectUserID(c *gin.Context, userID int64) ([]byte, string, error) {
	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/form-data") {
		return injectUserIDMultipart(c, userID)
	}
	return injectUserIDJSON(c, userID)
}

func injectUserIDJSON(c *gin.Context, userID int64) ([]byte, string, error) {
	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, maxAdBody))
	if err != nil {
		return nil, "", exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err)
	}

	payload, err := decodePayload(raw)
	if err != nil {
		return nil, "", err
	}
	payload["user_id"] = userID

	out, err := json.Marshal(payload)
	if err != nil {
		return nil, "", exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err)
	}
	return out, "application/json", nil
}

func injectUserIDMultipart(c *gin.Context, userID int64) ([]byte, string, error) {
	if err := c.Request.ParseMultipartForm(maxAdBody); err != nil {
		return nil, "", exception.NewAppError("INVALID_REQUEST", http.StatusBadRequest).WithCause(err)
	}

	payloadRaw := c.PostForm("payload")
	payload, err := decodePayload([]byte(payloadRaw))
	if err != nil {
		return nil, "", err
	}
	payload["user_id"] = userID

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
	delete(payload, "user_id")
	return payload, nil
}

func escapeQuotes(s string) string {
	if !utf8.ValidString(s) {
		s = strings.ToValidUTF8(s, "")
	}
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
