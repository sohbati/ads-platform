package impl

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"ads-platform/internal/business/ads/errorcode"
	"ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/ads/repository"
	"ads-platform/internal/business/ads/service"
	searchclient "ads-platform/internal/business/search/client"
	"ads-platform/internal/core/exception"
	"ads-platform/internal/core/storage"
)

const (
	maxTitleLen       = 120
	maxDescriptionLen = 10000
	maxFilenameLen    = 255
)

var allowedPriceTypes = map[string]struct{}{
	"fixed":      {},
	"negotiable": {},
	"free":       {},
	"salary":     {},
}

var extByContentType = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type adService struct {
	ads         repository.AdRepository
	images      repository.AdImageRepository
	catalog     searchclient.CatalogClient
	objects     storage.ObjectStorage
	maxPics     int
	maxPicBytes int64
}

func NewAdService(
	ads repository.AdRepository,
	images repository.AdImageRepository,
	catalog searchclient.CatalogClient,
	objects storage.ObjectStorage,
	maxPics int,
	maxPicBytes int64,
) service.AdService {
	if maxPics < 1 {
		maxPics = 8
	}
	if maxPicBytes < 1 {
		maxPicBytes = 10 << 20
	}
	return &adService{
		ads:         ads,
		images:      images,
		catalog:     catalog,
		objects:     objects,
		maxPics:     maxPics,
		maxPicBytes: maxPicBytes,
	}
}

func (s *adService) Create(ctx context.Context, in service.CreateAdInput) (*model.Ad, error) {
	if err := s.validate(in); err != nil {
		return nil, err
	}
	if err := s.resolveCatalog(ctx, in.CategoryID, in.CityID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	ad := &model.Ad{
		UserID:      in.UserID,
		CategoryID:  in.CategoryID,
		CityID:      in.CityID,
		Title:       strings.TrimSpace(in.Title),
		Description: strings.TrimSpace(in.Description),
		Status:      model.AdStatusActive,
		PriceAmount: in.PriceAmount,
		PriceType:   normalizePriceType(in.PriceType),
		Currency:    normalizeCurrency(in.Currency),
		Attrs:       normalizeJSONObject(in.Attrs),
		Media:       json.RawMessage("[]"),
		Contact:     normalizeJSONObject(in.Contact),
		Location:    buildLocation(in.Latitude, in.Longitude, in.Neighborhood),
		PublishedAt: &now,
	}

	if err := s.ads.Create(ctx, ad); err != nil {
		return nil, err
	}

	if len(in.Pictures) == 0 {
		return ad, nil
	}

	media, err := s.storePictures(ctx, ad, in.Pictures)
	if err != nil {
		return nil, err
	}
	raw, err := json.Marshal(media)
	if err != nil {
		return nil, err
	}
	if err := s.ads.UpdateMedia(ctx, ad.ID, raw); err != nil {
		return nil, err
	}
	ad.Media = raw
	return ad, nil
}

func (s *adService) validate(in service.CreateAdInput) error {
	if in.UserID <= 0 {
		return exception.NewAppError(errorcode.ErrAdInvalidUser.Code, errorcode.ErrAdInvalidUser.HttpStatus)
	}

	title := strings.TrimSpace(in.Title)
	if title == "" || len([]rune(title)) > maxTitleLen {
		return exception.NewAppError(errorcode.ErrAdInvalidTitle.Code, errorcode.ErrAdInvalidTitle.HttpStatus)
	}

	desc := strings.TrimSpace(in.Description)
	if desc == "" || len([]rune(desc)) > maxDescriptionLen {
		return exception.NewAppError(errorcode.ErrAdInvalidDescription.Code, errorcode.ErrAdInvalidDescription.HttpStatus)
	}

	if in.CategoryID <= 0 {
		return exception.NewAppError(errorcode.ErrAdInvalidCategory.Code, errorcode.ErrAdInvalidCategory.HttpStatus)
	}
	if in.CityID <= 0 {
		return exception.NewAppError(errorcode.ErrAdInvalidCity.Code, errorcode.ErrAdInvalidCity.HttpStatus)
	}

	if (in.Latitude == nil) != (in.Longitude == nil) {
		return exception.NewAppError(errorcode.ErrAdInvalidLocation.Code, errorcode.ErrAdInvalidLocation.HttpStatus)
	}
	if in.Latitude != nil {
		if *in.Latitude < -90 || *in.Latitude > 90 || *in.Longitude < -180 || *in.Longitude > 180 {
			return exception.NewAppError(errorcode.ErrAdInvalidLocation.Code, errorcode.ErrAdInvalidLocation.HttpStatus)
		}
	}

	if in.PriceAmount != nil && *in.PriceAmount < 0 {
		return exception.NewAppError(errorcode.ErrAdInvalidPrice.Code, errorcode.ErrAdInvalidPrice.HttpStatus)
	}
	if in.PriceType != "" {
		if _, ok := allowedPriceTypes[strings.ToLower(strings.TrimSpace(in.PriceType))]; !ok {
			return exception.NewAppError(errorcode.ErrAdInvalidPrice.Code, errorcode.ErrAdInvalidPrice.HttpStatus, in.PriceType)
		}
	}

	if len(in.Attrs) > 0 && !json.Valid(in.Attrs) {
		return exception.NewAppError(errorcode.ErrAdInvalidAttrs.Code, errorcode.ErrAdInvalidAttrs.HttpStatus)
	}
	if len(in.Contact) > 0 && !json.Valid(in.Contact) {
		return exception.NewAppError(errorcode.ErrAdInvalidAttrs.Code, errorcode.ErrAdInvalidAttrs.HttpStatus)
	}

	if len(in.Pictures) > s.maxPics {
		return exception.NewAppError(
			errorcode.ErrAdTooManyPictures.Code, errorcode.ErrAdTooManyPictures.HttpStatus, fmt.Sprintf("%d", s.maxPics))
	}
	for _, pic := range in.Pictures {
		if err := s.validatePicture(pic); err != nil {
			return err
		}
	}
	if len(in.Pictures) > 0 && s.objects == nil {
		return exception.NewAppError(errorcode.ErrAdStorageUnavailable.Code, errorcode.ErrAdStorageUnavailable.HttpStatus)
	}
	return nil
}

func (s *adService) validatePicture(pic service.PictureInput) error {
	name := strings.TrimSpace(pic.Filename)
	if name == "" || len(name) > maxFilenameLen {
		return exception.NewAppError(errorcode.ErrAdInvalidPicture.Code, errorcode.ErrAdInvalidPicture.HttpStatus, name)
	}
	if _, ok := extByContentType[normalizeContentType(pic.ContentType, name)]; !ok {
		return exception.NewAppError(errorcode.ErrAdInvalidPicture.Code, errorcode.ErrAdInvalidPicture.HttpStatus, pic.ContentType)
	}
	if pic.Size <= 0 {
		return exception.NewAppError(errorcode.ErrAdInvalidPicture.Code, errorcode.ErrAdInvalidPicture.HttpStatus)
	}
	if pic.Size > s.maxPicBytes {
		return exception.NewAppError(errorcode.ErrAdPictureTooLarge.Code, errorcode.ErrAdPictureTooLarge.HttpStatus)
	}
	return nil
}

func (s *adService) resolveCatalog(ctx context.Context, categoryID, cityID int) error {
	cats, err := s.catalog.CategoriesByIDs(ctx, []int{categoryID}, true)
	if err != nil {
		return exception.NewAppError(
			errorcode.ErrAdCatalogUnavailable.Code, errorcode.ErrAdCatalogUnavailable.HttpStatus).WithCause(err)
	}
	if len(cats) == 0 {
		return exception.NewAppError(errorcode.ErrAdInvalidCategory.Code, errorcode.ErrAdInvalidCategory.HttpStatus)
	}
	// Parent categories include descendants besides themselves; ads belong on leaves.
	if ids := cats[0].DescendantIDs; len(ids) > 1 {
		return exception.NewAppError(errorcode.ErrAdCategoryNotLeaf.Code, errorcode.ErrAdCategoryNotLeaf.HttpStatus)
	}

	cities, err := s.catalog.CitiesByIDs(ctx, []int{cityID})
	if err != nil {
		return exception.NewAppError(
			errorcode.ErrAdCatalogUnavailable.Code, errorcode.ErrAdCatalogUnavailable.HttpStatus).WithCause(err)
	}
	if len(cities) == 0 {
		return exception.NewAppError(errorcode.ErrAdInvalidCity.Code, errorcode.ErrAdInvalidCity.HttpStatus)
	}
	return nil
}

type mediaItem struct {
	ObjectKey   string `json:"object_key"`
	URL         string `json:"url"`
	Thumb       string `json:"thumb"`
	ContentType string `json:"content_type"`
	IsCover     bool   `json:"is_cover"`
}

func (s *adService) storePictures(ctx context.Context, ad *model.Ad, pics []service.PictureInput) ([]mediaItem, error) {
	out := make([]mediaItem, 0, len(pics))
	now := time.Now().UTC()

	for i, pic := range pics {
		contentType := normalizeContentType(pic.ContentType, pic.Filename)
		ext := extByContentType[contentType]
		key, err := buildObjectKey(ad.UserID, ad.ID, ext)
		if err != nil {
			return nil, err
		}

		hasher := sha256.New()
		reader := io.TeeReader(pic.Body, hasher)
		url, err := s.objects.Put(ctx, key, contentType, reader, pic.Size)
		if err != nil {
			return nil, exception.NewAppError(
				errorcode.ErrAdStorageUnavailable.Code, errorcode.ErrAdStorageUnavailable.HttpStatus).WithCause(err)
		}
		sum := hex.EncodeToString(hasher.Sum(nil))
		adID := ad.ID
		image := &model.AdImage{
			UserID:           ad.UserID,
			AdID:             &adID,
			ObjectKey:        key,
			OriginalFilename: strings.TrimSpace(pic.Filename),
			ContentType:      contentType,
			FileSize:         pic.Size,
			Status:           model.ImageStatusUploaded,
			Checksum:         &sum,
			UploadedAt:       &now,
		}
		if err := s.images.Create(ctx, image); err != nil {
			return nil, err
		}

		out = append(out, mediaItem{
			ObjectKey:   key,
			URL:         url,
			Thumb:       url,
			ContentType: contentType,
			IsCover:     i == 0,
		})
	}
	return out, nil
}

func buildLocation(lat, lng *float64, neighborhood string) json.RawMessage {
	loc := map[string]any{}
	if lat != nil && lng != nil {
		loc["lat"] = *lat
		loc["lng"] = *lng
	}
	if n := strings.TrimSpace(neighborhood); n != "" {
		loc["neighborhood"] = n
	}
	if len(loc) == 0 {
		return json.RawMessage("{}")
	}
	raw, err := json.Marshal(loc)
	if err != nil {
		return json.RawMessage("{}")
	}
	return raw
}

func normalizeJSONObject(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return json.RawMessage("{}")
	}
	return raw
}

func normalizePriceType(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	if _, ok := allowedPriceTypes[v]; ok {
		return v
	}
	return "fixed"
}

func normalizeCurrency(v string) string {
	v = strings.ToUpper(strings.TrimSpace(v))
	if len(v) != 3 {
		return "IRR"
	}
	return v
}

func normalizeContentType(ct, filename string) string {
	ct = strings.ToLower(strings.TrimSpace(ct))
	if _, ok := extByContentType[ct]; ok {
		return ct
	}
	switch strings.ToLower(path.Ext(filename)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	default:
		return ct
	}
}

func buildObjectKey(userID, adID int64, ext string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("ad picture: generate object key: %w", err)
	}
	return fmt.Sprintf("ads/%d/%d/%s%s", userID, adID, hex.EncodeToString(buf), ext), nil
}
