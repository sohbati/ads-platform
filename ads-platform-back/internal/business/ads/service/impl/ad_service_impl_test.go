package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"testing"

	"ads-platform/internal/business/ads/model"
	"ads-platform/internal/business/ads/service"
	searchclient "ads-platform/internal/business/search/client"
	"ads-platform/internal/core/exception"

	"gorm.io/gorm"
)

type fakeAdRepo struct {
	nextID int64
	ads    map[int64]model.Ad
}

func newFakeAdRepo() *fakeAdRepo {
	return &fakeAdRepo{nextID: 1, ads: map[int64]model.Ad{}}
}

func (f *fakeAdRepo) Create(_ context.Context, ad *model.Ad) error {
	ad.ID = f.nextID
	f.nextID++
	f.ads[ad.ID] = *ad
	return nil
}

func (f *fakeAdRepo) GetByID(_ context.Context, id int64) (*model.Ad, error) {
	ad, ok := f.ads[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copy := ad
	return &copy, nil
}

func (f *fakeAdRepo) Update(_ context.Context, ad *model.Ad) error {
	if _, ok := f.ads[ad.ID]; !ok {
		return gorm.ErrRecordNotFound
	}
	f.ads[ad.ID] = *ad
	return nil
}

func (f *fakeAdRepo) UpdateMedia(_ context.Context, id int64, media json.RawMessage) error {
	ad, ok := f.ads[id]
	if !ok {
		return errors.New("not found")
	}
	ad.Media = media
	f.ads[id] = ad
	return nil
}

func (f *fakeAdRepo) ListByUserID(_ context.Context, userID int64, limit int) ([]model.Ad, error) {
	out := make([]model.Ad, 0)
	for _, ad := range f.ads {
		if ad.UserID != userID || ad.Status == model.AdStatusDeleted {
			continue
		}
		out = append(out, ad)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

type fakeImageRepo struct {
	nextID int64
	rows   []model.AdImage
}

func (f *fakeImageRepo) Create(_ context.Context, image *model.AdImage) error {
	if f.nextID == 0 {
		f.nextID = 1
	}
	image.ID = f.nextID
	f.nextID++
	f.rows = append(f.rows, *image)
	return nil
}

func (f *fakeImageRepo) GetByID(context.Context, int64) (*model.AdImage, error) {
	return nil, errors.New("unused")
}

func (f *fakeImageRepo) Update(context.Context, *model.AdImage) error {
	return errors.New("unused")
}

func (f *fakeImageRepo) NextObjectSeq(_ context.Context, userID int64) (int64, error) {
	var n int64
	for i := range f.rows {
		if f.rows[i].UserID == userID {
			n++
		}
	}
	return n + 1, nil
}

type fakeCatalog struct {
	categories    []searchclient.Category
	categoriesErr error
	cities        []searchclient.City
	citiesErr     error
}

func (f *fakeCatalog) CategoriesBySlugs(context.Context, []string, bool) ([]searchclient.Category, error) {
	return f.categories, f.categoriesErr
}
func (f *fakeCatalog) CategoriesByIDs(context.Context, []int, bool) ([]searchclient.Category, error) {
	return f.categories, f.categoriesErr
}
func (f *fakeCatalog) CitiesBySlugs(context.Context, []string) ([]searchclient.City, error) {
	return f.cities, f.citiesErr
}
func (f *fakeCatalog) CitiesByIDs(context.Context, []int) ([]searchclient.City, error) {
	return f.cities, f.citiesErr
}
func (f *fakeCatalog) AttrSchemasByNames(context.Context, []string) ([]searchclient.AttrSchema, error) {
	return nil, nil
}

type fakeStorage struct {
	keys []string
	err  error
}

func (f *fakeStorage) Put(_ context.Context, key, _ string, body io.Reader, _ int64) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	_, _ = io.Copy(io.Discard, body)
	f.keys = append(f.keys, key)
	return "/ads-media/" + key, nil
}

func assertAppErrorCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error %s, got nil", code)
	}
	appErr, ok := exception.AsAppError(err)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.ErrorCode != code {
		t.Fatalf("expected %s, got %s", code, appErr.ErrorCode)
	}
}

func validInput() service.CreateAdInput {
	attrs := json.RawMessage(`{"rooms":2,"area":80}`)
	lat, lng := 35.6892, 51.3890
	price := int64(1_500_000)
	return service.CreateAdInput{
		UserID:       7,
		CategoryID:   13,
		CityID:       1,
		Title:        "Two-bedroom apartment",
		Description:  "Bright unit\nnear the metro.",
		Latitude:     &lat,
		Longitude:    &lng,
		Neighborhood: "Vanak",
		PriceAmount:  &price,
		PriceType:    "fixed",
		Currency:     "IRR",
		Attrs:        attrs,
	}
}

func leafCatalog() *fakeCatalog {
	return &fakeCatalog{
		categories: []searchclient.Category{{ID: 13, Slug: "apartment-sell", Title: "Apartment", DescendantIDs: []int{13}}},
		cities:     []searchclient.City{{ID: 1, Slug: "tehran", Name: "Tehran"}},
	}
}

func jpegPic(name string) service.PictureInput {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.NRGBA{R: 200, A: 255})
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
	return service.PictureInput{
		Filename:    name,
		ContentType: "image/jpeg",
		Size:        int64(buf.Len()),
		Body:        bytes.NewReader(buf.Bytes()),
	}
}

func pngPic(name string) service.PictureInput {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.NRGBA{G: 180, A: 255})
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return service.PictureInput{
		Filename:    name,
		ContentType: "image/png",
		Size:        int64(buf.Len()),
		Body:        bytes.NewReader(buf.Bytes()),
	}
}

func TestCreateAdHappyPath(t *testing.T) {
	ads := newFakeAdRepo()
	images := &fakeImageRepo{}
	store := &fakeStorage{}
	svc := NewAdService(ads, images, leafCatalog(), store, 8, 10<<20)

	in := validInput()
	in.Pictures = []service.PictureInput{jpegPic("cover.jpg")}

	ad, err := svc.Create(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ad.ID == 0 || ad.Status != model.AdStatusActive || ad.PublishedAt == nil {
		t.Errorf("ad not persisted as active: %+v", ad)
	}
	if !bytes.Contains(ad.Location, []byte(`"lat":35.6892`)) || !bytes.Contains(ad.Location, []byte("Vanak")) {
		t.Errorf("location = %s", ad.Location)
	}
	if !bytes.Contains(ad.Attrs, []byte(`"rooms":2`)) {
		t.Errorf("attrs = %s", ad.Attrs)
	}
	if len(images.rows) != 1 || images.rows[0].AdID == nil || *images.rows[0].AdID != ad.ID {
		t.Fatalf("image row not linked to ad: %+v", images.rows)
	}
	if len(store.keys) != 2 || store.keys[0] != "ads/7/7_1.webp" || store.keys[1] != "ads/7/7_1-t.webp" {
		t.Errorf("object keys = %v", store.keys)
	}
	if images.rows[0].ContentType != "image/webp" {
		t.Errorf("stored content type = %s", images.rows[0].ContentType)
	}
	if !bytes.Contains(ad.Media, []byte(`"is_cover":true`)) || !bytes.Contains(ad.Media, []byte("/ads-media/ads/7/7_1-t.webp")) {
		t.Errorf("media = %s", ad.Media)
	}
	if bytes.Contains(ad.Media, []byte("http://")) || bytes.Contains(ad.Media, []byte("https://")) {
		t.Errorf("media URLs should be relative: %s", ad.Media)
	}
}

func TestCreateAdPictureKeysFollowUserSequence(t *testing.T) {
	ads := newFakeAdRepo()
	images := &fakeImageRepo{}
	store := &fakeStorage{}
	svc := NewAdService(ads, images, leafCatalog(), store, 8, 10<<20)

	first := validInput()
	first.Pictures = []service.PictureInput{jpegPic("a.jpg"), pngPic("b.png")}
	if _, err := svc.Create(context.Background(), first); err != nil {
		t.Fatalf("first ad: %v", err)
	}

	second := validInput()
	second.Title = "Second listing"
	second.Pictures = []service.PictureInput{jpegPic("c.jpg")}
	if _, err := svc.Create(context.Background(), second); err != nil {
		t.Fatalf("second ad: %v", err)
	}

	other := validInput()
	other.UserID = 99
	other.Pictures = []service.PictureInput{jpegPic("d.jpg")}
	if _, err := svc.Create(context.Background(), other); err != nil {
		t.Fatalf("other user: %v", err)
	}

	want := []string{
		"ads/7/7_1.webp", "ads/7/7_1-t.webp",
		"ads/7/7_2.webp", "ads/7/7_2-t.webp",
		"ads/7/7_3.webp", "ads/7/7_3-t.webp",
		"ads/99/99_1.webp", "ads/99/99_1-t.webp",
	}
	if len(store.keys) != len(want) {
		t.Fatalf("keys=%v, want %v", store.keys, want)
	}
	for i, key := range want {
		if store.keys[i] != key {
			t.Fatalf("keys=%v, want %v", store.keys, want)
		}
	}
}

func TestCreateAdRejectsCorruptPicture(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, leafCatalog(), &fakeStorage{}, 8, 10<<20)
	in := validInput()
	in.Pictures = []service.PictureInput{{
		Filename: "cover.jpg", ContentType: "image/jpeg", Size: 4, Body: bytes.NewReader([]byte("nope")),
	}}
	_, err := svc.Create(context.Background(), in)
	assertAppErrorCode(t, err, "AD_INVALID_PICTURE")
}

func TestCreateAdWithoutPictures(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	ad, err := svc.Create(context.Background(), validInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(ad.Media) != "[]" {
		t.Errorf("media = %s, want []", ad.Media)
	}
}

func TestCreateAdValidation(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, leafCatalog(), &fakeStorage{}, 2, 100)

	cases := []struct {
		name     string
		mutate   func(*service.CreateAdInput)
		wantCode string
	}{
		{"missing user", func(in *service.CreateAdInput) { in.UserID = 0 }, "AD_INVALID_USER"},
		{"empty title", func(in *service.CreateAdInput) { in.Title = "  " }, "AD_INVALID_TITLE"},
		{"empty description", func(in *service.CreateAdInput) { in.Description = "" }, "AD_INVALID_DESCRIPTION"},
		{"lat without lng", func(in *service.CreateAdInput) { in.Longitude = nil }, "AD_INVALID_LOCATION"},
		{"bad coords", func(in *service.CreateAdInput) { lat := 200.0; in.Latitude = &lat }, "AD_INVALID_LOCATION"},
		{"negative price", func(in *service.CreateAdInput) { v := int64(-1); in.PriceAmount = &v }, "AD_INVALID_PRICE"},
		{"bad price type", func(in *service.CreateAdInput) { in.PriceType = "barter" }, "AD_INVALID_PRICE"},
		{"bad attrs json", func(in *service.CreateAdInput) { in.Attrs = json.RawMessage(`{`) }, "AD_INVALID_ATTRS"},
		{"too many pictures", func(in *service.CreateAdInput) {
			in.Pictures = []service.PictureInput{
				{Filename: "a.jpg", ContentType: "image/jpeg", Size: 1, Body: bytes.NewReader([]byte("a"))},
				{Filename: "b.jpg", ContentType: "image/jpeg", Size: 1, Body: bytes.NewReader([]byte("b"))},
				{Filename: "c.jpg", ContentType: "image/jpeg", Size: 1, Body: bytes.NewReader([]byte("c"))},
			}
		}, "AD_TOO_MANY_PICTURES"},
		{"picture too large", func(in *service.CreateAdInput) {
			in.Pictures = []service.PictureInput{{Filename: "a.jpg", ContentType: "image/jpeg", Size: 101, Body: bytes.NewReader(nil)}}
		}, "AD_PICTURE_TOO_LARGE"},
		{"svg rejected", func(in *service.CreateAdInput) {
			in.Pictures = []service.PictureInput{{Filename: "a.svg", ContentType: "image/svg+xml", Size: 10, Body: bytes.NewReader(nil)}}
		}, "AD_INVALID_PICTURE"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := validInput()
			tc.mutate(&in)
			_, err := svc.Create(context.Background(), in)
			assertAppErrorCode(t, err, tc.wantCode)
		})
	}
}

func TestCreateAdRejectsParentCategory(t *testing.T) {
	catalog := &fakeCatalog{
		categories: []searchclient.Category{{ID: 5, Slug: "vehicles", DescendantIDs: []int{5, 7, 9}}},
		cities:     []searchclient.City{{ID: 1}},
	}
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, catalog, nil, 8, 10<<20)
	_, err := svc.Create(context.Background(), validInput())
	assertAppErrorCode(t, err, "AD_CATEGORY_NOT_LEAF")
}

func TestCreateAdUnknownCategory(t *testing.T) {
	catalog := &fakeCatalog{
		categories: nil,
		cities:     []searchclient.City{{ID: 1}},
	}
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, catalog, nil, 8, 10<<20)
	_, err := svc.Create(context.Background(), validInput())
	assertAppErrorCode(t, err, "AD_INVALID_CATEGORY")
}

func TestCreateAdPicturesNeedStorage(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	in := validInput()
	in.Pictures = []service.PictureInput{{
		Filename: "a.jpg", ContentType: "image/jpeg", Size: 1, Body: bytes.NewReader([]byte("a")),
	}}
	_, err := svc.Create(context.Background(), in)
	assertAppErrorCode(t, err, "AD_STORAGE_UNAVAILABLE")
}

func TestListByUserReturnsOnlyOwnerAds(t *testing.T) {
	ads := newFakeAdRepo()
	svc := NewAdService(ads, &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)

	mine, err := svc.Create(context.Background(), validInput())
	if err != nil {
		t.Fatal(err)
	}
	other := validInput()
	other.UserID = 99
	if _, err := svc.Create(context.Background(), other); err != nil {
		t.Fatal(err)
	}
	deleted := ads.ads[mine.ID]
	deleted.Status = model.AdStatusDeleted
	ads.ads[mine.ID] = deleted

	again := validInput()
	again.Title = "Second listing"
	live, err := svc.Create(context.Background(), again)
	if err != nil {
		t.Fatal(err)
	}

	items, err := svc.ListByUser(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].ID != live.ID {
		t.Fatalf("items=%+v", items)
	}
	if items[0].CityName != "Tehran" || items[0].Neighborhood != "Vanak" {
		t.Fatalf("projection=%+v", items[0])
	}
}

func TestListByUserRejectsInvalidUser(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	_, err := svc.ListByUser(context.Background(), 0)
	assertAppErrorCode(t, err, "AD_INVALID_USER")
}

func TestCreateAdPreservesMultilineDescription(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	in := validInput()
	in.Description = "line one\n\nline three"
	ad, err := svc.Create(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	if ad.Description != "line one\n\nline three" {
		t.Errorf("description = %q", ad.Description)
	}
}

func TestGetPublicReturnsActiveWithFullMedia(t *testing.T) {
	ads := newFakeAdRepo()
	svc := NewAdService(ads, &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	ad, err := svc.Create(context.Background(), validInput())
	if err != nil {
		t.Fatal(err)
	}
	ad.Media = json.RawMessage(`[{"url":"/ads-media/ads/7/7_1.webp","thumb":"/ads-media/ads/7/7_1-t.webp","is_cover":true},{"url":"/ads-media/ads/7/7_2.webp","thumb":"/ads-media/ads/7/7_2-t.webp"}]`)
	ads.ads[ad.ID] = *ad

	got, err := svc.GetPublic(context.Background(), ad.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != ad.Title || got.Description != ad.Description || got.CityName != "Tehran" || got.Neighborhood != "Vanak" {
		t.Fatalf("public ad: %+v", got)
	}
	if len(got.Media) != 2 || got.Media[0].URL != "/ads-media/ads/7/7_1.webp" || got.Media[1].URL != "/ads-media/ads/7/7_2.webp" {
		t.Fatalf("media: %+v", got.Media)
	}

	_, err = svc.GetPublic(context.Background(), 0)
	assertAppErrorCode(t, err, "AD_NOT_FOUND")
	_, err = svc.GetPublic(context.Background(), 99)
	assertAppErrorCode(t, err, "AD_NOT_FOUND")

	deleted := ads.ads[ad.ID]
	deleted.Status = model.AdStatusDeleted
	ads.ads[ad.ID] = deleted
	_, err = svc.GetPublic(context.Background(), ad.ID)
	assertAppErrorCode(t, err, "AD_NOT_FOUND")
}

func TestGetForOwnerHidesOthersAndDeleted(t *testing.T) {
	ads := newFakeAdRepo()
	svc := NewAdService(ads, &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	ad, err := svc.Create(context.Background(), validInput())
	if err != nil {
		t.Fatal(err)
	}

	got, err := svc.GetForOwner(context.Background(), 7, ad.ID)
	if err != nil || got.ID != ad.ID {
		t.Fatalf("owner get: %+v %v", got, err)
	}
	_, err = svc.GetForOwner(context.Background(), 8, ad.ID)
	assertAppErrorCode(t, err, "AD_NOT_FOUND")

	deleted := ads.ads[ad.ID]
	deleted.Status = model.AdStatusDeleted
	ads.ads[ad.ID] = deleted
	_, err = svc.GetForOwner(context.Background(), 7, ad.ID)
	assertAppErrorCode(t, err, "AD_NOT_FOUND")
}

func TestUpdateAdOwnerFields(t *testing.T) {
	ads := newFakeAdRepo()
	svc := NewAdService(ads, &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	ad, err := svc.Create(context.Background(), validInput())
	if err != nil {
		t.Fatal(err)
	}

	in := validInput()
	in.Title = "Edited title"
	in.Description = "Edited body"
	in.Neighborhood = "Jordan"
	in.Latitude = nil
	in.Longitude = nil
	price := int64(2000)
	in.PriceAmount = &price
	updated, err := svc.Update(context.Background(), ad.ID, in)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Title != "Edited title" || updated.Description != "Edited body" {
		t.Fatalf("updated=%+v", updated)
	}
	if !bytes.Contains(updated.Location, []byte("Jordan")) {
		t.Fatalf("location=%s", updated.Location)
	}
	if !bytes.Contains(updated.Location, []byte("35.6892")) {
		t.Fatalf("expected preserved coords, location=%s", updated.Location)
	}
}

func TestUpdateKeepsAndDropsExistingMedia(t *testing.T) {
	ads := newFakeAdRepo()
	svc := NewAdService(ads, &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	ad, err := svc.Create(context.Background(), validInput())
	if err != nil {
		t.Fatal(err)
	}
	ad.Media = json.RawMessage(`[{"url":"/ads-media/ads/7/7_1.webp","thumb":"/ads-media/ads/7/7_1-t.webp","is_cover":true},{"url":"/ads-media/ads/7/7_2.webp","thumb":"/ads-media/ads/7/7_2-t.webp"}]`)
	ads.ads[ad.ID] = *ad

	keep := []string{"/ads-media/ads/7/7_2.webp"}
	in := validInput()
	in.KeepMedia = &keep
	updated, err := svc.Update(context.Background(), ad.ID, in)
	if err != nil {
		t.Fatal(err)
	}
	var media []map[string]any
	if err := json.Unmarshal(updated.Media, &media); err != nil {
		t.Fatal(err)
	}
	if len(media) != 1 || media[0]["url"] != "/ads-media/ads/7/7_2.webp" {
		t.Fatalf("media=%s", updated.Media)
	}
	if media[0]["is_cover"] != true {
		t.Fatalf("expected remaining photo to be cover, media=%s", updated.Media)
	}

	empty := []string{}
	in.KeepMedia = &empty
	updated, err = svc.Update(context.Background(), ad.ID, in)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(updated.Media, &media); err != nil {
		t.Fatal(err)
	}
	if len(media) != 0 {
		t.Fatalf("expected no media, got %s", updated.Media)
	}
}

func TestUpdateIgnoresKeepMediaFromOtherAd(t *testing.T) {
	ads := newFakeAdRepo()
	svc := NewAdService(ads, &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	ad, err := svc.Create(context.Background(), validInput())
	if err != nil {
		t.Fatal(err)
	}
	ad.Media = json.RawMessage(`[{"url":"/ads-media/ads/7/7_1.webp","is_cover":true}]`)
	ads.ads[ad.ID] = *ad

	keep := []string{"/ads-media/ads/99/99_1.webp"}
	in := validInput()
	in.KeepMedia = &keep
	updated, err := svc.Update(context.Background(), ad.ID, in)
	if err != nil {
		t.Fatal(err)
	}
	var media []map[string]any
	if err := json.Unmarshal(updated.Media, &media); err != nil {
		t.Fatal(err)
	}
	if len(media) != 0 {
		t.Fatalf("foreign url should be dropped, media=%s", updated.Media)
	}
}

func TestUpdateAdRejectsOtherUser(t *testing.T) {
	ads := newFakeAdRepo()
	svc := NewAdService(ads, &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	ad, err := svc.Create(context.Background(), validInput())
	if err != nil {
		t.Fatal(err)
	}
	in := validInput()
	in.UserID = 99
	in.Title = "Hack"
	_, err = svc.Update(context.Background(), ad.ID, in)
	assertAppErrorCode(t, err, "AD_NOT_FOUND")
}
