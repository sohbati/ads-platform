package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"ads-platform-ui/internal/business/newad/viewmodel"
	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"

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
		Locale:         string(loc),
		CityID:         vm.CityID,
		CitySlug:       pageData.CitySlug,
		CityName:       i18n.CityDisplayName(h.i18n, h.cities, loc, pageData.CitySlug),
		MaxPictures:    maxAdPictures,
		SuccessHref:    "/my-info/user-ads",
		Categories:     categories,
		Schemas:        schemas,
		Enums:          enums,
	}

	c.HTML(http.StatusOK, "new_ad", vm)
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
