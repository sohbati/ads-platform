package handler

import (
	"encoding/json"
	"io"
	"net/http"
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
