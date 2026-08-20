package impl

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	adsmodel "ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/search/client"
	"ads-platform/internal/business/search/errorcode"
	"ads-platform/internal/business/search/model"
	"ads-platform/internal/business/search/repository"
	"ads-platform/internal/business/search/service"
	"ads-platform/internal/core/exception"
)

type searchService struct {
	repo    repository.SearchRepository
	catalog client.CatalogClient
}

func NewSearchService(repo repository.SearchRepository, catalog client.CatalogClient) service.SearchService {
	return &searchService{repo: repo, catalog: catalog}
}

func (s *searchService) Search(
	ctx context.Context,
	placeSlug, categorySlug string,
	citiesCSV string,
	query string,
	priceMin, priceMax *int64,
	hasPhoto *bool,
	neighborhood, sort string,
	page, limit int,
) (*model.SearchResponse, error) {
	placeSlug = strings.TrimSpace(placeSlug)
	categorySlug = strings.TrimSpace(categorySlug)

	cat, categoryIDs, err := s.resolveCategory(ctx, categorySlug)
	if err != nil {
		return nil, err
	}

	cityIDs, err := s.resolveCityIDs(ctx, placeSlug, citiesCSV)
	if err != nil {
		return nil, err
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 24
	}
	if sort == "" {
		sort = "newest"
	}

	filter := model.SearchFilter{
		PlaceSlug:    placeSlug,
		CategorySlug: categorySlug,
		CategoryIDs:  categoryIDs,
		CityIDs:      cityIDs,
		Query:        strings.TrimSpace(query),
		PriceMin:     priceMin,
		PriceMax:     priceMax,
		HasPhoto:     hasPhoto,
		Neighborhood: strings.TrimSpace(neighborhood),
		Sort:         sort,
		Page:         page,
		Limit:        limit,
	}

	ads, total, err := s.repo.SearchAds(ctx, filter)
	if err != nil {
		return nil, err
	}

	cityNames := s.cityNamesFor(ctx, ads)

	items := make([]model.SearchAdItem, 0, len(ads))
	for _, ad := range ads {
		items = append(items, toSearchItem(ad, cityNames))
	}

	return &model.SearchResponse{
		Place:         placeSlug,
		Category:      categorySlug,
		CategoryTitle: cat.Title,
		Filters: model.SearchFilters{
			Cities:       cityIDs,
			Query:        filter.Query,
			PriceMin:     priceMin,
			PriceMax:     priceMax,
			HasPhoto:     hasPhoto,
			Neighborhood: filter.Neighborhood,
			Sort:         sort,
		},
		Pagination: model.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
		Ads: items,
	}, nil
}

// resolveCategory resolves the slug via cache-service and returns the category
// plus self+descendant ids so parent categories include their leaves.
func (s *searchService) resolveCategory(ctx context.Context, slug string) (client.Category, []int, error) {
	slug = strings.ToLower(strings.TrimSpace(slug))
	if slug == "" {
		return client.Category{}, nil, exception.NewAppError(
			errorcode.ErrInvalidCategory.Code, errorcode.ErrInvalidCategory.HttpStatus, slug)
	}

	categories, err := s.catalog.CategoriesBySlugs(ctx, []string{slug}, true)
	if err != nil {
		return client.Category{}, nil, exception.NewAppError(
			errorcode.ErrCatalogUnavailable.Code, errorcode.ErrCatalogUnavailable.HttpStatus).WithCause(err)
	}
	if len(categories) == 0 {
		return client.Category{}, nil, exception.NewAppError(
			errorcode.ErrInvalidCategory.Code, errorcode.ErrInvalidCategory.HttpStatus, slug)
	}

	cat := categories[0]
	ids := cat.DescendantIDs
	if len(ids) == 0 {
		ids = []int{cat.ID}
	}
	return cat, ids, nil
}

func (s *searchService) resolveCityIDs(ctx context.Context, placeSlug, citiesCSV string) ([]int, error) {
	if isIranPlace(placeSlug) {
		ids, err := parseCityIDs(citiesCSV)
		if err != nil {
			return nil, err
		}
		if len(ids) == 0 {
			// iran without cities = nationwide (no city filter)
			return nil, nil
		}
		if err := s.validateCityIDs(ctx, ids); err != nil {
			return nil, err
		}
		return ids, nil
	}

	cities, err := s.catalog.CitiesBySlugs(ctx, []string{placeSlug})
	if err != nil {
		return nil, exception.NewAppError(
			errorcode.ErrCatalogUnavailable.Code, errorcode.ErrCatalogUnavailable.HttpStatus).WithCause(err)
	}
	if len(cities) == 0 {
		return nil, exception.NewAppError(errorcode.ErrInvalidPlace.Code, errorcode.ErrInvalidPlace.HttpStatus, placeSlug)
	}
	return []int{cities[0].ID}, nil
}

func (s *searchService) validateCityIDs(ctx context.Context, ids []int) error {
	cities, err := s.catalog.CitiesByIDs(ctx, ids)
	if err != nil {
		return exception.NewAppError(
			errorcode.ErrCatalogUnavailable.Code, errorcode.ErrCatalogUnavailable.HttpStatus).WithCause(err)
	}

	known := make(map[int]struct{}, len(cities))
	for _, city := range cities {
		known[city.ID] = struct{}{}
	}
	for _, id := range ids {
		if _, ok := known[id]; !ok {
			return exception.NewAppError(
				errorcode.ErrInvalidCities.Code, errorcode.ErrInvalidCities.HttpStatus, strconv.Itoa(id))
		}
	}
	return nil
}

// cityNamesFor batch-resolves city names for the result page; enrichment is
// best-effort, so lookup failures leave names empty rather than failing search.
func (s *searchService) cityNamesFor(ctx context.Context, ads []adsmodel.Ad) map[int]string {
	seen := map[int]struct{}{}
	ids := make([]int, 0, len(ads))
	for _, ad := range ads {
		if _, ok := seen[ad.CityID]; ok {
			continue
		}
		seen[ad.CityID] = struct{}{}
		ids = append(ids, ad.CityID)
	}
	if len(ids) == 0 {
		return nil
	}

	cities, err := s.catalog.CitiesByIDs(ctx, ids)
	if err != nil {
		return nil
	}

	names := make(map[int]string, len(cities))
	for _, city := range cities {
		names[city.ID] = city.Name
	}
	return names
}

func isIranPlace(slug string) bool {
	return strings.EqualFold(strings.TrimSpace(slug), "iran")
}

func parseCityIDs(csv string) ([]int, error) {
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return nil, nil
	}
	parts := strings.Split(csv, ",")
	ids := make([]int, 0, len(parts))
	seen := map[int]struct{}{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.Atoi(p)
		if err != nil || id <= 0 {
			return nil, exception.NewAppError(errorcode.ErrInvalidCities.Code, errorcode.ErrInvalidCities.HttpStatus, csv)
		}
		if _, dup := seen[id]; dup {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

func toSearchItem(ad adsmodel.Ad, cityNames map[int]string) model.SearchAdItem {
	item := model.SearchAdItem{
		ID:          ad.ID,
		Title:       ad.Title,
		PriceAmount: ad.PriceAmount,
		PriceType:   ad.PriceType,
		Currency:    ad.Currency,
		CityID:      ad.CityID,
		CategoryID:  ad.CategoryID,
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
		} else if url, ok := media[0]["url"].(string); ok {
			item.Thumbnail = url
		}
	}

	var loc map[string]any
	if len(ad.Location) > 0 && json.Unmarshal(ad.Location, &loc) == nil {
		if n, ok := loc["neighborhood"].(string); ok {
			item.Neighborhood = n
		}
	}
	return item
}
