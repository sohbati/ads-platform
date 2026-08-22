package handler

import (
	"net/http"
	"strconv"

	"ads-platform/internal/business/ads/errorcode"
	"ads-platform/internal/business/ads/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

// AdImageHandler handles HTTP requests for ad image metadata.
// The caller (ads-bff) is trusted to supply the authenticated user_id.
type AdImageHandler struct {
	imageService service.AdImageService
}

func NewAdImageHandler(imageService service.AdImageService) *AdImageHandler {
	return &AdImageHandler{imageService: imageService}
}

type registerImageRequest struct {
	UserID           int64  `json:"user_id"`
	OriginalFilename string `json:"original_filename"`
	ContentType      string `json:"content_type"`
	FileSize         int64  `json:"file_size"`
}

// Register handles POST /api/v1/ads/images
func (h *AdImageHandler) Register(c *gin.Context) {
	var req registerImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(
			errorcode.ErrImageInvalidSize.Code, http.StatusBadRequest).WithCause(err), 0)
		return
	}

	image, err := h.imageService.Register(c.Request.Context(), req.UserID, req.OriginalFilename, req.ContentType, req.FileSize)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusCreated, image)
}

type markUploadedRequest struct {
	UserID   int64  `json:"user_id"`
	Checksum string `json:"checksum"`
	FileSize int64  `json:"file_size"`
}

// MarkUploaded handles POST /api/v1/ads/images/:id/uploaded
func (h *AdImageHandler) MarkUploaded(c *gin.Context) {
	imageID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	var req markUploadedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.HandleError(c, exception.NewAppError(
			errorcode.ErrImageChecksumRequired.Code, http.StatusBadRequest).WithCause(err), 0)
		return
	}

	image, err := h.imageService.MarkUploaded(c.Request.Context(), req.UserID, imageID, req.Checksum, req.FileSize)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusOK, image)
}

// Get handles GET /api/v1/ads/images/:id?user_id=
func (h *AdImageHandler) Get(c *gin.Context) {
	imageID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	image, err := h.imageService.Get(c.Request.Context(), userID, imageID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusOK, image)
}

// Delete handles DELETE /api/v1/ads/images/:id?user_id=
func (h *AdImageHandler) Delete(c *gin.Context) {
	imageID, err := parseID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	if err := h.imageService.Delete(c.Request.Context(), userID, imageID); err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exception.NewAppError(
			errorcode.ErrImageNotFound.Code, errorcode.ErrImageNotFound.HttpStatus, raw)
	}
	return id, nil
}
