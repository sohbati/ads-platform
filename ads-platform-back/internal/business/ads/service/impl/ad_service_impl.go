package impl

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	"ads-platform/internal/core/imageconv"
	"ads-platform/internal/core/storage"

	"gorm.io/gorm"
)

const (
	maxTitleLen       = 120
	maxDescriptionLen = 10000
	maxFilenameLen    = 255
	maxUserAdsList    = 200
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

func (s *adService) GetForOwner(ctx context.Context, userID, adID int64) (*model.Ad, error) {
	return s.loadOwned(ctx, userID, adID)
}

func (s *adService) GetPublic(ctx context.Context, adID int64) (*model.PublicAd, error) {
	if adID <= 0 {
		return nil, exception.NewAppError(errorcode.ErrAdNotFound.Code, errorcode.ErrAdNotFound.HttpStatus)
	}
	ad, err := s.ads.GetByID(ctx, adID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(errorcode.ErrAdNotFound.Code, errorcode.ErrAdNotFound.HttpStatus)
		}
		return nil, err
	}
	if ad == nil || ad.Status != model.AdStatusActive {
		return nil, exception.NewAppError(errorcode.ErrAdNotFound.Code, errorcode.ErrAdNotFound.HttpStatus)
	}
	return toPublicAd(ad, s.cityNames(ctx, []model.Ad{*ad})), nil
}

func (s *adService) Update(ctx context.Context, adID int64, in service.CreateAdInput) (*model.Ad, error) {
	existing, err := s.loadOwned(ctx, in.UserID, adID)
	if err != nil {
		return nil, err
	}
	if err := s.validate(in); err != nil {
		return nil, err
	}
	if err := s.resolveCatalog(ctx, in.CategoryID, in.CityID); err != nil {
		return nil, err
	}

	lat, lng := in.Latitude, in.Longitude
	if lat == nil || lng == nil {
		existLat, existLng := coordsFromLocation(existing.Location)
		if lat == nil {
			lat = existLat
		}
		if lng == nil {
			lng = existLng
		}
	}

	existing.CategoryID = in.CategoryID
	existing.CityID = in.CityID
	existing.Title = strings.TrimSpace(in.Title)
	existing.Description = strings.TrimSpace(in.Description)
	existing.PriceAmount = in.PriceAmount
	existing.PriceType = normalizePriceType(in.PriceType)
	existing.Currency = normalizeCurrency(in.Currency)
	existing.Attrs = normalizeJSONObject(in.Attrs)
	if len(in.Contact) > 0 {
		existing.Contact = normalizeJSONObject(in.Contact)
	}
	existing.Location = buildLocation(lat, lng, in.Neighborhood)

	if len(in.Pictures) > 0 {
		media, err := s.storePictures(ctx, existing, in.Pictures)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(media)
		if err != nil {
			return nil, err
		}
		existing.Media = raw
	}

	if err := s.ads.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *adService) loadOwned(ctx context.Context, userID, adID int64) (*model.Ad, error) {
	if userID <= 0 {
		return nil, exception.NewAppError(errorcode.ErrAdInvalidUser.Code, errorcode.ErrAdInvalidUser.HttpStatus)
	}
	if adID <= 0 {
		return nil, exception.NewAppError(errorcode.ErrAdNotFound.Code, errorcode.ErrAdNotFound.HttpStatus)
	}

	ad, err := s.ads.GetByID(ctx, adID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewAppError(errorcode.ErrAdNotFound.Code, errorcode.ErrAdNotFound.HttpStatus)
		}
		return nil, err
	}
	if ad == nil || ad.UserID != userID || ad.Status == model.AdStatusDeleted {
		return nil, exception.NewAppError(errorcode.ErrAdNotFound.Code, errorcode.ErrAdNotFound.HttpStatus)
	}
	return ad, nil
}

func coordsFromLocation(raw json.RawMessage) (*float64, *float64) {
	var loc map[string]any
	if len(raw) == 0 || json.Unmarshal(raw, &loc) != nil {
		return nil, nil
	}
	lat, okLat := jsonNumber(loc["lat"])
	lng, okLng := jsonNumber(loc["lng"])
	if !okLat || !okLng {
		return nil, nil
	}
	return &lat, &lng
}

func jsonNumber(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func (s *adService) ListByUser(ctx context.Context, userID int64) ([]model.UserAdItem, error) {
	if userID <= 0 {
		return nil, exception.NewAppError(errorcode.ErrAdInvalidUser.Code, errorcode.ErrAdInvalidUser.HttpStatus)
	}

	ads, err := s.ads.ListByUserID(ctx, userID, maxUserAdsList)
	if err != nil {
		return nil, err
	}

	cityNames := s.cityNames(ctx, ads)
	items := make([]model.UserAdItem, 0, len(ads))
	for _, ad := range ads {
		items = append(items, toUserAdItem(ad, cityNames))
	}
	return items, nil
}

func (s *adService) cityNames(ctx context.Context, ads []model.Ad) map[int]string {
	ids := uniqueCityIDs(ads)
	if len(ids) == 0 || s.catalog == nil {
		return map[int]string{}
	}
	cities, err := s.catalog.CitiesByIDs(ctx, ids)
	if err != nil {
		return map[int]string{}
	}
	out := make(map[int]string, len(cities))
	for _, city := range cities {
		out[city.ID] = city.Name
	}
	return out
}

func uniqueCityIDs(ads []model.Ad) []int {
	seen := make(map[int]struct{}, len(ads))
	ids := make([]int, 0, len(ads))
	for _, ad := range ads {
		if ad.CityID <= 0 {
			continue
		}
		if _, ok := seen[ad.CityID]; ok {
			continue
		}
		seen[ad.CityID] = struct{}{}
		ids = append(ids, ad.CityID)
	}
	return ids
}

func toUserAdItem(ad model.Ad, cityNames map[int]string) model.UserAdItem {
	item := model.UserAdItem{
		ID:          ad.ID,
		Title:       ad.Title,
		PriceAmount: ad.PriceAmount,
		PriceType:   ad.PriceType,
		Currency:    ad.Currency,
		CityID:      ad.CityID,
		CategoryID:  ad.CategoryID,
		Status:      ad.Status,
		HasPhoto:    false,
	}
	if name, ok := cityNames[ad.CityID]; ok {
		item.CityName = name
	}
	if ad.Slug != nil {
		item.Slug = *ad.Slug
	}
	if ad.PublishedAt != nil {
		s := ad.PublishedAt.UTC().Format(time.RFC3339)
		item.PublishedAt = &s
	}

	var media []map[string]any
	if len(ad.Media) > 0 && json.Unmarshal(ad.Media, &media) == nil && len(media) > 0 {
		item.HasPhoto = true
		if thumb, ok := media[0]["thumb"].(string); ok && thumb != "" {
			item.Thumbnail = thumb
		} else if u, ok := media[0]["url"].(string); ok {
			item.Thumbnail = u
		}
	}

	item.Neighborhood = neighborhoodFromLocation(ad.Location)
	return item
}

func toPublicAd(ad *model.Ad, cityNames map[int]string) *model.PublicAd {
	out := &model.PublicAd{
		ID:          ad.ID,
		Title:       ad.Title,
		Description: ad.Description,
		PriceAmount: ad.PriceAmount,
		PriceType:   ad.PriceType,
		Currency:    ad.Currency,
		CityID:      ad.CityID,
		Media:       publicMedia(ad.Media),
	}
	if name, ok := cityNames[ad.CityID]; ok {
		out.CityName = name
	}
	out.Neighborhood = neighborhoodFromLocation(ad.Location)
	if ad.PublishedAt != nil {
		s := ad.PublishedAt.UTC().Format(time.RFC3339)
		out.PublishedAt = &s
	}
	return out
}

func neighborhoodFromLocation(raw json.RawMessage) string {
	var loc map[string]any
	if len(raw) > 0 && json.Unmarshal(raw, &loc) == nil {
		if n, ok := loc["neighborhood"].(string); ok {
			return n
		}
	}
	return ""
}

func publicMedia(raw json.RawMessage) []model.PublicMedia {
	var items []map[string]any
	if len(raw) == 0 || json.Unmarshal(raw, &items) != nil || len(items) == 0 {
		return []model.PublicMedia{}
	}
	out := make([]model.PublicMedia, 0, len(items))
	for _, item := range items {
		url, _ := item["url"].(string)
		if url == "" {
			continue
		}
		m := model.PublicMedia{URL: url}
		if thumb, ok := item["thumb"].(string); ok {
			m.Thumb = thumb
		}
		if cover, ok := item["is_cover"].(bool); ok {
			m.IsCover = cover
		}
		out = append(out, m)
	}
	return out
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

	seq, err := s.images.NextObjectSeq(ctx, ad.UserID)
	if err != nil {
		return nil, err
	}

	for i, pic := range pics {
		contentType := normalizeContentType(pic.ContentType, pic.Filename)
		raw, err := io.ReadAll(io.LimitReader(pic.Body, s.maxPicBytes+1))
		if err != nil {
			return nil, exception.NewAppError(
				errorcode.ErrAdInvalidPicture.Code, errorcode.ErrAdInvalidPicture.HttpStatus).WithCause(err)
		}
		variants, err := imageconv.FromUpload(bytes.NewReader(raw), contentType)
		if err != nil {
			return nil, exception.NewAppError(
				errorcode.ErrAdInvalidPicture.Code, errorcode.ErrAdInvalidPicture.HttpStatus).WithCause(err)
		}

		n := seq + int64(i)
		fullKey := buildObjectKey(ad.UserID, n, ".webp")
		thumbKey := buildObjectKey(ad.UserID, n, "-t.webp")

		url, err := s.objects.Put(ctx, fullKey, "image/webp", bytes.NewReader(variants.Full), int64(len(variants.Full)))
		if err != nil {
			return nil, exception.NewAppError(
				errorcode.ErrAdStorageUnavailable.Code, errorcode.ErrAdStorageUnavailable.HttpStatus).WithCause(err)
		}
		thumbURL, err := s.objects.Put(ctx, thumbKey, "image/webp", bytes.NewReader(variants.Thumb), int64(len(variants.Thumb)))
		if err != nil {
			return nil, exception.NewAppError(
				errorcode.ErrAdStorageUnavailable.Code, errorcode.ErrAdStorageUnavailable.HttpStatus).WithCause(err)
		}

		sum := sha256.Sum256(variants.Full)
		checksum := hex.EncodeToString(sum[:])
		adID := ad.ID
		image := &model.AdImage{
			UserID:           ad.UserID,
			AdID:             &adID,
			ObjectKey:        fullKey,
			OriginalFilename: strings.TrimSpace(pic.Filename),
			ContentType:      "image/webp",
			FileSize:         int64(len(variants.Full)),
			Status:           model.ImageStatusUploaded,
			Checksum:         &checksum,
			UploadedAt:       &now,
		}
		if err := s.images.Create(ctx, image); err != nil {
			return nil, err
		}

		out = append(out, mediaItem{
			ObjectKey:   fullKey,
			URL:         url,
			Thumb:       thumbURL,
			ContentType: "image/webp",
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

func buildObjectKey(userID, seq int64, ext string) string {
	return fmt.Sprintf("ads/%d/%d_%d%s", userID, userID, seq, ext)
}
