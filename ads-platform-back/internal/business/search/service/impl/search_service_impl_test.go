package impl

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	adsmodel "ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/search/client"
	"ads-platform/internal/business/search/errorcode"
	"ads-platform/internal/business/search/model"
	"ads-platform/internal/core/exception"
)

// fakeRepo records the filter it was called with and returns canned results.
type fakeRepo struct {
	gotFilter model.SearchFilter
	ads       []adsmodel.Ad
	total     int64
	err       error
}

func (f *fakeRepo) SearchAds(_ context.Context, filter model.SearchFilter) ([]adsmodel.Ad, int64, error) {
	f.gotFilter = filter
	return f.ads, f.total, f.err
}

// fakeCatalog returns canned catalog responses.
type fakeCatalog struct {
	categories       []client.Category
	categoriesErr    error
	citiesBySlugs    []client.City
	citiesBySlugsErr error
	citiesByIDs      []client.City
	citiesByIDsErr   error
}

func (f *fakeCatalog) CategoriesBySlugs(context.Context, []string, bool) ([]client.Category, error) {
	return f.categories, f.categoriesErr
}

func (f *fakeCatalog) CitiesBySlugs(context.Context, []string) ([]client.City, error) {
	return f.citiesBySlugs, f.citiesBySlugsErr
}

func (f *fakeCatalog) CitiesByIDs(context.Context, []int) ([]client.City, error) {
	return f.citiesByIDs, f.citiesByIDsErr
}

func assertAppErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %s, got nil", code)
	}
	appErr, ok := exception.AsAppError(err)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.ErrorCode != code {
		t.Fatalf("expected error code %s, got %s", code, appErr.ErrorCode)
	}
}

func i64(v int64) *int64 { return &v }

func TestSearchHappyPath(t *testing.T) {
	published := time.Date(2026, 8, 1, 10, 30, 0, 0, time.UTC)
	slug := "used-bike"
	repo := &fakeRepo{
		ads: []adsmodel.Ad{
			{
				ID:          42,
				CategoryID:  7,
				CityID:      1,
				Title:       "Bike for sale",
				PriceAmount: i64(500000),
				PriceType:   "fixed",
				Currency:    "IRR",
				Media:       json.RawMessage(`[{"thumb":"t.jpg","url":"u.jpg"}]`),
				Location:    json.RawMessage(`{"neighborhood":"Vanak"}`),
				Slug:        &slug,
				PublishedAt: &published,
			},
		},
		total: 1,
	}
	catalog := &fakeCatalog{
		categories:  []client.Category{{ID: 5, Slug: "vehicles", Title: "Vehicles", DescendantIDs: []int{5, 7}}},
		citiesByIDs: []client.City{{ID: 1, Slug: "tehran", Name: "Tehran"}},
	}
	svc := NewSearchService(repo, catalog)

	resp, err := svc.Search(context.Background(), "iran", "vehicles", "1", " bike ", i64(100), i64(900), nil, "", "", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Filter passed to the repository.
	f := repo.gotFilter
	if got, want := f.CategoryIDs, []int{5, 7}; len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("CategoryIDs = %v, want %v", got, want)
	}
	if len(f.CityIDs) != 1 || f.CityIDs[0] != 1 {
		t.Errorf("CityIDs = %v, want [1]", f.CityIDs)
	}
	if f.Query != "bike" {
		t.Errorf("Query = %q, want %q (trimmed)", f.Query, "bike")
	}
	if f.Page != 1 || f.Limit != 24 {
		t.Errorf("Page/Limit = %d/%d, want defaults 1/24", f.Page, f.Limit)
	}
	if f.Sort != "newest" {
		t.Errorf("Sort = %q, want default %q", f.Sort, "newest")
	}

	// Response envelope.
	if resp.CategoryTitle != "Vehicles" {
		t.Errorf("CategoryTitle = %q, want %q", resp.CategoryTitle, "Vehicles")
	}
	if resp.Pagination.Total != 1 || resp.Pagination.Page != 1 || resp.Pagination.Limit != 24 {
		t.Errorf("Pagination = %+v, want total 1, page 1, limit 24", resp.Pagination)
	}
	if len(resp.Ads) != 1 {
		t.Fatalf("expected 1 ad, got %d", len(resp.Ads))
	}

	// Ad projection including city name enrichment and JSONB flattening.
	item := resp.Ads[0]
	if item.ID != 42 || item.Title != "Bike for sale" {
		t.Errorf("item identity = %d/%q, want 42/\"Bike for sale\"", item.ID, item.Title)
	}
	if item.CityName != "Tehran" {
		t.Errorf("CityName = %q, want %q", item.CityName, "Tehran")
	}
	if !item.HasPhoto || item.Thumbnail != "t.jpg" {
		t.Errorf("HasPhoto/Thumbnail = %v/%q, want true/%q", item.HasPhoto, item.Thumbnail, "t.jpg")
	}
	if item.Neighborhood != "Vanak" {
		t.Errorf("Neighborhood = %q, want %q", item.Neighborhood, "Vanak")
	}
	if item.PublishedAt == nil || *item.PublishedAt != "2026-08-01T10:30:00Z" {
		t.Errorf("PublishedAt = %v, want 2026-08-01T10:30:00Z", item.PublishedAt)
	}
}

func TestSearchIranWithoutCitiesIsNationwide(t *testing.T) {
	repo := &fakeRepo{}
	catalog := &fakeCatalog{
		categories: []client.Category{{ID: 5, Slug: "vehicles", Title: "Vehicles"}},
	}
	svc := NewSearchService(repo, catalog)

	if _, err := svc.Search(context.Background(), "iran", "vehicles", "", "", nil, nil, nil, "", "", 1, 24); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.gotFilter.CityIDs != nil {
		t.Errorf("CityIDs = %v, want nil (nationwide)", repo.gotFilter.CityIDs)
	}
}

func TestSearchCityPlaceSlugResolvesToCityID(t *testing.T) {
	repo := &fakeRepo{}
	catalog := &fakeCatalog{
		categories:    []client.Category{{ID: 5, Slug: "vehicles"}},
		citiesBySlugs: []client.City{{ID: 3, Slug: "shiraz", Name: "Shiraz"}},
	}
	svc := NewSearchService(repo, catalog)

	if _, err := svc.Search(context.Background(), "shiraz", "vehicles", "", "", nil, nil, nil, "", "", 1, 24); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.gotFilter.CityIDs) != 1 || repo.gotFilter.CityIDs[0] != 3 {
		t.Errorf("CityIDs = %v, want [3]", repo.gotFilter.CityIDs)
	}
}

func TestSearchAllCategoriesSkipsCategoryFilter(t *testing.T) {
	repo := &fakeRepo{}
	// Category lookup erroring proves "all" never hits the catalog.
	catalog := &fakeCatalog{categoriesErr: errors.New("must not be called")}
	svc := NewSearchService(repo, catalog)

	resp, err := svc.Search(context.Background(), "iran", "all", "", "phone", nil, nil, nil, "", "", 1, 24)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.gotFilter.CategoryIDs != nil {
		t.Errorf("CategoryIDs = %v, want nil (no category filter)", repo.gotFilter.CategoryIDs)
	}
	if resp.Category != "all" || resp.CategoryTitle != "" {
		t.Errorf("Category/CategoryTitle = %q/%q, want \"all\"/empty", resp.Category, resp.CategoryTitle)
	}
}

func TestSearchCategoryWithoutDescendantsUsesOwnID(t *testing.T) {
	repo := &fakeRepo{}
	catalog := &fakeCatalog{
		categories: []client.Category{{ID: 9, Slug: "books"}},
	}
	svc := NewSearchService(repo, catalog)

	if _, err := svc.Search(context.Background(), "iran", "books", "", "", nil, nil, nil, "", "", 1, 24); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.gotFilter.CategoryIDs) != 1 || repo.gotFilter.CategoryIDs[0] != 9 {
		t.Errorf("CategoryIDs = %v, want [9]", repo.gotFilter.CategoryIDs)
	}
}

func TestSearchErrors(t *testing.T) {
	validCategory := []client.Category{{ID: 5, Slug: "vehicles"}}

	cases := []struct {
		name     string
		place    string
		category string
		cities   string
		catalog  *fakeCatalog
		wantCode string
	}{
		{
			name:     "empty category slug",
			place:    "iran",
			category: "  ",
			catalog:  &fakeCatalog{},
			wantCode: errorcode.ErrInvalidCategory.Code,
		},
		{
			name:     "unknown category slug",
			place:    "iran",
			category: "nope",
			catalog:  &fakeCatalog{categories: nil},
			wantCode: errorcode.ErrInvalidCategory.Code,
		},
		{
			name:     "catalog down on category lookup",
			place:    "iran",
			category: "vehicles",
			catalog:  &fakeCatalog{categoriesErr: errors.New("boom")},
			wantCode: errorcode.ErrCatalogUnavailable.Code,
		},
		{
			name:     "unknown place slug",
			place:    "atlantis",
			category: "vehicles",
			catalog:  &fakeCatalog{categories: validCategory, citiesBySlugs: nil},
			wantCode: errorcode.ErrInvalidPlace.Code,
		},
		{
			name:     "catalog down on place lookup",
			place:    "tehran",
			category: "vehicles",
			catalog:  &fakeCatalog{categories: validCategory, citiesBySlugsErr: errors.New("boom")},
			wantCode: errorcode.ErrCatalogUnavailable.Code,
		},
		{
			name:     "malformed cities csv",
			place:    "iran",
			category: "vehicles",
			cities:   "1,abc",
			catalog:  &fakeCatalog{categories: validCategory},
			wantCode: errorcode.ErrInvalidCities.Code,
		},
		{
			name:     "unknown city id",
			place:    "iran",
			category: "vehicles",
			cities:   "999",
			catalog:  &fakeCatalog{categories: validCategory, citiesByIDs: nil},
			wantCode: errorcode.ErrInvalidCities.Code,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewSearchService(&fakeRepo{}, tc.catalog)
			_, err := svc.Search(context.Background(), tc.place, tc.category, tc.cities, "", nil, nil, nil, "", "", 1, 24)
			assertAppErrorCode(t, err, tc.wantCode)
		})
	}
}

func TestSearchRepositoryErrorPropagates(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &fakeRepo{err: repoErr}
	catalog := &fakeCatalog{categories: []client.Category{{ID: 5, Slug: "vehicles"}}}
	svc := NewSearchService(repo, catalog)

	_, err := svc.Search(context.Background(), "iran", "vehicles", "", "", nil, nil, nil, "", "", 1, 24)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error to propagate, got %v", err)
	}
}

func TestSearchPaginationClamping(t *testing.T) {
	cases := []struct {
		name              string
		page, limit       int
		wantPage, wantLim int
	}{
		{"zero values get defaults", 0, 0, 1, 24},
		{"negative page becomes 1", -3, 10, 1, 10},
		{"limit above 50 resets to 24", 2, 100, 2, 24},
		{"valid values pass through", 3, 50, 3, 50},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeRepo{}
			catalog := &fakeCatalog{categories: []client.Category{{ID: 5, Slug: "vehicles"}}}
			svc := NewSearchService(repo, catalog)

			if _, err := svc.Search(context.Background(), "iran", "vehicles", "", "", nil, nil, nil, "", "", tc.page, tc.limit); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo.gotFilter.Page != tc.wantPage || repo.gotFilter.Limit != tc.wantLim {
				t.Errorf("Page/Limit = %d/%d, want %d/%d",
					repo.gotFilter.Page, repo.gotFilter.Limit, tc.wantPage, tc.wantLim)
			}
		})
	}
}

func TestCityNamesLookupFailureIsBestEffort(t *testing.T) {
	// City-name enrichment failing must not fail the search itself.
	repo := &fakeRepo{
		ads:   []adsmodel.Ad{{ID: 1, CityID: 4, Title: "Chair"}},
		total: 1,
	}
	catalog := &fakeCatalog{
		categories:     []client.Category{{ID: 5, Slug: "home"}},
		citiesByIDsErr: errors.New("boom"),
	}
	svc := NewSearchService(repo, catalog)

	resp, err := svc.Search(context.Background(), "iran", "home", "", "", nil, nil, nil, "", "", 1, 24)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Ads[0].CityName != "" {
		t.Errorf("CityName = %q, want empty when lookup fails", resp.Ads[0].CityName)
	}
}

func TestToSearchItemMediaFallbacks(t *testing.T) {
	base := adsmodel.Ad{ID: 1, Title: "x"}

	noMedia := base
	item := toSearchItem(noMedia, nil)
	if item.HasPhoto || item.Thumbnail != "" {
		t.Errorf("no media: HasPhoto/Thumbnail = %v/%q, want false/empty", item.HasPhoto, item.Thumbnail)
	}

	urlOnly := base
	urlOnly.Media = json.RawMessage(`[{"url":"full.jpg"}]`)
	item = toSearchItem(urlOnly, nil)
	if !item.HasPhoto || item.Thumbnail != "full.jpg" {
		t.Errorf("url fallback: HasPhoto/Thumbnail = %v/%q, want true/full.jpg", item.HasPhoto, item.Thumbnail)
	}

	emptyArray := base
	emptyArray.Media = json.RawMessage(`[]`)
	item = toSearchItem(emptyArray, nil)
	if item.HasPhoto {
		t.Error("empty media array: HasPhoto = true, want false")
	}

	badJSON := base
	badJSON.Media = json.RawMessage(`{not json`)
	item = toSearchItem(badJSON, nil)
	if item.HasPhoto {
		t.Error("malformed media: HasPhoto = true, want false")
	}
}

func TestParseCityIDs(t *testing.T) {
	valid := []struct {
		in   string
		want []int
	}{
		{"", nil},
		{"  ", nil},
		{"1,2,3", []int{1, 2, 3}},
		{" 1 , 2 ", []int{1, 2}},
		{"1,1,2", []int{1, 2}}, // duplicates removed
		{"1,,2", []int{1, 2}},  // empty parts skipped
	}
	for _, tc := range valid {
		got, err := parseCityIDs(tc.in)
		if err != nil {
			t.Fatalf("parseCityIDs(%q): unexpected error: %v", tc.in, err)
		}
		if len(got) != len(tc.want) {
			t.Fatalf("parseCityIDs(%q) = %v, want %v", tc.in, got, tc.want)
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Fatalf("parseCityIDs(%q) = %v, want %v", tc.in, got, tc.want)
			}
		}
	}

	for _, in := range []string{"abc", "1,x", "0", "-5"} {
		if _, err := parseCityIDs(in); err == nil {
			t.Errorf("parseCityIDs(%q): expected error, got nil", in)
		}
	}
}

func TestIsIranPlace(t *testing.T) {
	for _, s := range []string{"iran", "IRAN", " Iran "} {
		if !isIranPlace(s) {
			t.Errorf("isIranPlace(%q) = false, want true", s)
		}
	}
	if isIranPlace("tehran") {
		t.Error("isIranPlace(\"tehran\") = true, want false")
	}
}
