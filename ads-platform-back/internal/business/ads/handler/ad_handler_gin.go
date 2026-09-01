package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"ads-platform/internal/business/ads/errorcode"
	"ads-platform/internal/business/ads/service"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/middleware"

	"github.com/gin-gonic/gin"
)

type AdHandler struct {
	ads service.AdService
}

func NewAdHandler(ads service.AdService) *AdHandler {
	return &AdHandler{ads: ads}
}

type createAdRequest struct {
	UserID       int64           `json:"user_id"`
	CategoryID   int             `json:"category_id"`
	CityID       int             `json:"city_id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Latitude     *float64        `json:"latitude"`
	Longitude    *float64        `json:"longitude"`
	Neighborhood string          `json:"neighborhood"`
	PriceAmount  *int64          `json:"price_amount"`
	PriceType    string          `json:"price_type"`
	Currency     string          `json:"currency"`
	Attrs        json.RawMessage `json:"attrs"`
	Contact      json.RawMessage `json:"contact"`
	KeepMedia    *[]string       `json:"keep_media"`
}

// Create handles POST /api/v1/ads
// JSON body for ads without pictures, or multipart with a "payload" JSON field
// and "pictures" files.
func (h *AdHandler) Create(c *gin.Context) {
	req, pics, err := parseCreateRequest(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	defer closePictures(pics)

	ad, err := h.ads.Create(c.Request.Context(), service.CreateAdInput{
		UserID:       req.UserID,
		CategoryID:   req.CategoryID,
		CityID:       req.CityID,
		Title:        req.Title,
		Description:  req.Description,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Neighborhood: req.Neighborhood,
		PriceAmount:  req.PriceAmount,
		PriceType:    req.PriceType,
		Currency:     req.Currency,
		Attrs:        req.Attrs,
		Contact:      req.Contact,
		Pictures:     pics,
	})
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	c.JSON(http.StatusCreated, ad)
}

// ListByUser handles GET /api/v1/users/:userId/ads
func (h *AdHandler) ListByUser(c *gin.Context) {
	userID, err := parseUserID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	items, err := h.ads.ListByUser(c.Request.Context(), userID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ads": items})
}

// ListStats handles GET /api/v1/users/:userId/ad-stats
func (h *AdHandler) ListStats(c *gin.Context) {
	userID, err := parseUserID(c.Param("userId"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	result, err := h.ads.ListStats(c.Request.Context(), userID, c.Query("from"), c.Query("to"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetPublic handles GET /api/v1/ads/:id
func (h *AdHandler) GetPublic(c *gin.Context) {
	adID, err := parseAdID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	ad, err := h.ads.GetPublic(c.Request.Context(), adID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, ad)
}

// GetPublicContact handles GET /api/v1/ads/:id/contact
func (h *AdHandler) GetPublicContact(c *gin.Context) {
	adID, err := parseAdID(c.Param("id"))
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	contact, err := h.ads.GetPublicContact(c.Request.Context(), adID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, contact)
}

// GetForOwner handles GET /api/v1/users/:userId/ads/:adId
func (h *AdHandler) GetForOwner(c *gin.Context) {
	userID, adID, err := parseOwnerIDs(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	ad, err := h.ads.GetForOwner(c.Request.Context(), userID, adID)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, ad)
}

// Update handles PUT /api/v1/users/:userId/ads/:adId
func (h *AdHandler) Update(c *gin.Context) {
	userID, adID, err := parseOwnerIDs(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}

	req, pics, err := parseCreateRequest(c)
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	defer closePictures(pics)

	ad, err := h.ads.Update(c.Request.Context(), adID, service.CreateAdInput{
		UserID:       userID,
		CategoryID:   req.CategoryID,
		CityID:       req.CityID,
		Title:        req.Title,
		Description:  req.Description,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Neighborhood: req.Neighborhood,
		PriceAmount:  req.PriceAmount,
		PriceType:    req.PriceType,
		Currency:     req.Currency,
		Attrs:        req.Attrs,
		Contact:      req.Contact,
		Pictures:     pics,
		KeepMedia:    req.KeepMedia,
	})
	if err != nil {
		middleware.HandleError(c, err, 0)
		return
	}
	c.JSON(http.StatusOK, ad)
}

func parseOwnerIDs(c *gin.Context) (int64, int64, error) {
	userID, err := parseUserID(c.Param("userId"))
	if err != nil {
		return 0, 0, err
	}
	adID, err := parseAdID(c.Param("adId"))
	if err != nil {
		return 0, 0, err
	}
	return userID, adID, nil
}

func parseAdID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exception.NewAppError(errorcode.ErrAdNotFound.Code, errorcode.ErrAdNotFound.HttpStatus, raw)
	}
	return id, nil
}

func parseUserID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, exception.NewAppError(errorcode.ErrAdInvalidUser.Code, errorcode.ErrAdInvalidUser.HttpStatus, raw)
	}
	return id, nil
}

func parseCreateRequest(c *gin.Context) (createAdRequest, []service.PictureInput, error) {
	ct := c.ContentType()
	if strings.HasPrefix(ct, "multipart/form-data") {
		return parseMultipart(c)
	}

	var req createAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, nil, exception.NewAppError(
			errorcode.ErrAdInvalidTitle.Code, http.StatusBadRequest).WithCause(err)
	}
	return req, nil, nil
}

func parseMultipart(c *gin.Context) (createAdRequest, []service.PictureInput, error) {
	var req createAdRequest
	payload := c.PostForm("payload")
	if payload == "" {
		return req, nil, exception.NewAppError(errorcode.ErrAdInvalidTitle.Code, http.StatusBadRequest)
	}
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		return req, nil, exception.NewAppError(
			errorcode.ErrAdInvalidTitle.Code, http.StatusBadRequest).WithCause(err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return req, nil, exception.NewAppError(
			errorcode.ErrAdInvalidPicture.Code, http.StatusBadRequest).WithCause(err)
	}

	files := form.File["pictures"]
	pics := make([]service.PictureInput, 0, len(files))
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			closePictures(pics)
			return req, nil, exception.NewAppError(
				errorcode.ErrAdInvalidPicture.Code, http.StatusBadRequest).WithCause(err)
		}
		pics = append(pics, service.PictureInput{
			Filename:    fh.Filename,
			ContentType: fh.Header.Get("Content-Type"),
			Size:        fh.Size,
			Body:        f,
		})
	}
	return req, pics, nil
}

func closePictures(pics []service.PictureInput) {
	for _, p := range pics {
		if c, ok := p.Body.(io.Closer); ok {
			_ = c.Close()
		}
	}
}
