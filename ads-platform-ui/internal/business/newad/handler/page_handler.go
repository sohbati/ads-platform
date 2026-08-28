package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"ads-platform-ui/internal/business/newad/viewmodel"
	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
	"ads-platform-ui/internal/core/media"

	"github.com/gin-gonic/gin"
)

const maxAdPictures = 8

type PageHandler struct {
	config *config.Config
	i18n   *i18n.Registry
	cities *cities.Catalog
	cdn    *cdn.Client
	bff    *bff.Client
}

func NewPageHandler(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, cdnClient *cdn.Client, bffClient *bff.Client) *PageHandler {
	return &PageHandler{config: cfg, i18n: reg, cities: catalog, cdn: cdnClient, bff: bffClient}
}

func (h *PageHandler) Index(c *gin.Context) {
	if _, ok := h.currentUser(c); !ok {
		target := "/query-ads?login=1&next=" + url.QueryEscape("/new-ad")
		c.Redirect(http.StatusFound, target)
		return
	}

	loc := i18n.FromContext(c)
	t := h.i18n.MessagesFor(loc)
	pageData := i18n.BuildPage(h.i18n, h.cities, loc, h.config.AppName, i18n.CityFromContext(c), c.Request.URL.Path, i18n.LocationsFromContext(c))
	pageData.DefaultCountryCode = h.config.DefaultCountryCode
	pageData.Title = h.config.AppName + " — " + t.Nav.NewAd
	pageData.Heading = t.Nav.NewAd
	pageData.IsAuthenticated = true

	vm := viewmodel.NewAdPage{Page: pageData}
	city, ok := h.cities.Get(pageData.CitySlug)
	if ok {
		vm.CityID = city.ID
	}

	categories, catErr := h.cdn.GetCategories(c.Request.Context())
	schemas, schemaErr := h.cdn.GetAttrSchemas(c.Request.Context())
	enums, enumErr := h.cdn.GetAttrEnums(c.Request.Context())
	if catErr != nil || schemaErr != nil {
		vm.LoadError = t.NewAd.LoadError
	}
	if enums == nil || enumErr != nil {
		enums = json.RawMessage("{}")
	}

	vm.Bootstrap = viewmodel.Bootstrap{
		Locale:      string(loc),
		CityID:      vm.CityID,
		CitySlug:    pageData.CitySlug,
		CityName:    i18n.CityDisplayName(h.i18n, h.cities, loc, pageData.CitySlug),
		MaxPictures: maxAdPictures,
		SuccessHref: "/my-info/user-ads",
		Categories:  categories,
		Schemas:     schemas,
		Enums:       enums,
	}

	c.HTML(http.StatusOK, "new_ad", vm)
}

func (h *PageHandler) Edit(c *gin.Context) {
	adID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || adID <= 0 {
		c.Redirect(http.StatusFound, "/my-info/user-ads")
		return
	}

	editPath := "/edit-ad/" + strconv.FormatInt(adID, 10)
	if _, ok := h.currentUser(c); !ok {
		target := "/query-ads?login=1&next=" + url.QueryEscape(editPath)
		c.Redirect(http.StatusFound, target)
		return
	}

	loc := i18n.FromContext(c)
	t := h.i18n.MessagesFor(loc)
	pageData := i18n.BuildPage(h.i18n, h.cities, loc, h.config.AppName, i18n.CityFromContext(c), c.Request.URL.Path, i18n.LocationsFromContext(c))
	pageData.DefaultCountryCode = h.config.DefaultCountryCode
	pageData.Title = h.config.AppName + " — " + t.NewAd.EditHeading
	pageData.Heading = t.NewAd.EditHeading
	pageData.IsAuthenticated = true

	vm := viewmodel.NewAdPage{Page: pageData}

	result, err := h.bff.Get(c.Request.Context(), "/api/v1/me/ads/"+strconv.FormatInt(adID, 10), bff.RequestCookies(c.Request))
	if err != nil || result.StatusCode == http.StatusUnauthorized {
		c.Redirect(http.StatusFound, "/query-ads?login=1&next="+url.QueryEscape(editPath))
		return
	}
	if result.StatusCode == http.StatusNotFound {
		vm.LoadError = t.NewAd.NotFound
		c.HTML(http.StatusNotFound, "new_ad", vm)
		return
	}
	if result.StatusCode != http.StatusOK {
		vm.LoadError = t.NewAd.LoadError
		c.HTML(http.StatusOK, "new_ad", vm)
		return
	}

	ad, err := parseAdJSON(result.Body)
	if err != nil {
		vm.LoadError = t.NewAd.LoadError
		c.HTML(http.StatusOK, "new_ad", vm)
		return
	}

	categories, catErr := h.cdn.GetCategories(c.Request.Context())
	schemas, schemaErr := h.cdn.GetAttrSchemas(c.Request.Context())
	enums, enumErr := h.cdn.GetAttrEnums(c.Request.Context())
	if catErr != nil || schemaErr != nil {
		vm.LoadError = t.NewAd.LoadError
		c.HTML(http.StatusOK, "new_ad", vm)
		return
	}
	if enums == nil || enumErr != nil {
		enums = json.RawMessage("{}")
	}

	cityID := ad.CityID
	citySlug := pageData.CitySlug
	cityName := i18n.CityDisplayName(h.i18n, h.cities, loc, citySlug)
	if rec, ok := h.cities.GetByID(cityID); ok {
		citySlug = rec.Slug
		cityName = i18n.CityDisplayName(h.i18n, h.cities, loc, rec.Slug)
	}
	vm.CityID = cityID

	vm.Bootstrap = viewmodel.Bootstrap{
		Locale:      string(loc),
		Mode:        "edit",
		AdID:        ad.ID,
		CityID:      cityID,
		CitySlug:    citySlug,
		CityName:    cityName,
		MaxPictures: maxAdPictures,
		SuccessHref: "/my-info/user-ads",
		Prefill:     adToPrefill(ad, h.config.MediaCDNURL),
		Categories:  categories,
		Schemas:     schemas,
		Enums:       enums,
	}

	c.HTML(http.StatusOK, "new_ad", vm)
}

type adJSON struct {
	ID          int64           `json:"id"`
	CategoryID  int             `json:"category_id"`
	CityID      int             `json:"city_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	PriceAmount *int64          `json:"price_amount"`
	PriceType   string          `json:"price_type"`
	Attrs       json.RawMessage `json:"attrs"`
	Media       json.RawMessage `json:"media"`
	Location    json.RawMessage `json:"location"`
}

func parseAdJSON(body []byte) (adJSON, error) {
	var ad adJSON
	err := json.Unmarshal(body, &ad)
	return ad, err
}

func adToPrefill(ad adJSON, mediaCDN string) *viewmodel.Prefill {
	prefill := &viewmodel.Prefill{
		CategoryID:   ad.CategoryID,
		Title:        ad.Title,
		Description:  ad.Description,
		PriceAmount:  ad.PriceAmount,
		PriceType:    ad.PriceType,
		Neighborhood: neighborhoodFromLocation(ad.Location),
		Attrs:        ad.Attrs,
		Media:        mediaFromJSON(ad.Media, mediaCDN),
	}
	return prefill
}

func neighborhoodFromLocation(raw json.RawMessage) string {
	var loc map[string]any
	if len(raw) == 0 || json.Unmarshal(raw, &loc) != nil {
		return ""
	}
	if n, ok := loc["neighborhood"].(string); ok {
		return n
	}
	return ""
}

func mediaFromJSON(raw json.RawMessage, mediaCDN string) []viewmodel.PrefillMedia {
	var items []map[string]any
	if len(raw) == 0 || json.Unmarshal(raw, &items) != nil {
		return nil
	}
	out := make([]viewmodel.PrefillMedia, 0, len(items))
	for _, item := range items {
		m := viewmodel.PrefillMedia{}
		if u, ok := item["url"].(string); ok {
			m.URL = media.PublicURL(mediaCDN, u)
		}
		if t, ok := item["thumb"].(string); ok {
			m.Thumb = media.PublicURL(mediaCDN, t)
		}
		if m.URL == "" && m.Thumb == "" {
			continue
		}
		out = append(out, m)
	}
	return out
}

func (h *PageHandler) currentUser(c *gin.Context) (*sessionUser, bool) {
	result, err := h.bff.Get(c.Request.Context(), "/api/v1/auth/me", bff.RequestCookies(c.Request))
	if err != nil || result.StatusCode != http.StatusOK {
		return nil, false
	}

	var resp meResponse
	if err := json.Unmarshal(result.Body, &resp); err != nil || !resp.Authenticated || resp.User == nil {
		return nil, false
	}
	return resp.User, true
}

type meResponse struct {
	Authenticated bool         `json:"authenticated"`
	User          *sessionUser `json:"user"`
}

type sessionUser struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
}
