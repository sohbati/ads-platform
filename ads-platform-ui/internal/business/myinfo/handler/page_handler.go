package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"ads-platform-ui/internal/business/myinfo/viewmodel"
	queryadsvm "ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
	"ads-platform-ui/internal/core/media"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	config *config.Config
	i18n   *i18n.Registry
	cities *cities.Catalog
	bff    *bff.Client
}

func NewPageHandler(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, client *bff.Client) *PageHandler {
	return &PageHandler{config: cfg, i18n: reg, cities: catalog, bff: client}
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

func (h *PageHandler) protectedPage(c *gin.Context, heading string) (i18n.Page, bool) {
	user, ok := h.currentUser(c)
	if !ok {
		target := "/query-ads?login=1&next=" + url.QueryEscape(c.Request.URL.Path)
		c.Redirect(http.StatusFound, target)
		return i18n.Page{}, false
	}

	pageData := i18n.BuildPage(h.i18n, h.cities, i18n.FromContext(c), h.config.AppName, i18n.CityFromContext(c), c.Request.URL.Path, i18n.LocationsFromContext(c))
	pageData.DefaultCountryCode = h.config.DefaultCountryCode
	pageData.Title = h.config.AppName + " — " + heading
	pageData.Heading = heading
	pageData.IsAuthenticated = true
	pageData.SessionUserName = user.Name
	pageData.SessionUserMobile = user.Mobile
	return pageData, true
}

func (h *PageHandler) renderProtected(c *gin.Context, templateName, heading string) {
	pageData, ok := h.protectedPage(c, heading)
	if !ok {
		return
	}
	c.HTML(http.StatusOK, templateName, pageData)
}

func (h *PageHandler) Index(c *gin.Context) {
	h.renderProtected(c, "myinfo", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.MyInfo)
}

func (h *PageHandler) UserDetails(c *gin.Context) {
	h.renderProtected(c, "myinfo_user_details", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.UserDetails)
}

func (h *PageHandler) UserAds(c *gin.Context) {
	t := h.i18n.MessagesFor(i18n.FromContext(c))
	pageData, ok := h.protectedPage(c, t.Nav.UserAds)
	if !ok {
		return
	}

	ads, unavailable := h.fetchUserAds(c, t)
	c.HTML(http.StatusOK, "myinfo_user_ads", viewmodel.UserAdsPage{
		Page:        pageData,
		Ads:         ads,
		Unavailable: unavailable,
	})
}

func (h *PageHandler) fetchUserAds(c *gin.Context, t i18n.Messages) ([]queryadsvm.SearchAd, bool) {
	result, err := h.bff.Get(c.Request.Context(), "/api/v1/me/ads", bff.RequestCookies(c.Request))
	if err != nil || result.StatusCode != http.StatusOK {
		return nil, true
	}

	var resp userAdsResponse
	if err := json.Unmarshal(result.Body, &resp); err != nil {
		return nil, true
	}

	out := make([]queryadsvm.SearchAd, 0, len(resp.Ads))
	for _, ad := range resp.Ads {
		out = append(out, toSearchAd(ad, t, h.config.MediaCDNURL))
	}
	return out, false
}

func toSearchAd(ad userAdJSON, t i18n.Messages, mediaCDN string) queryadsvm.SearchAd {
	out := queryadsvm.SearchAd{
		ID:        ad.ID,
		Title:     ad.Title,
		Location:  ad.CityName,
		Thumbnail: media.PublicURL(mediaCDN, ad.Thumbnail),
		HasPhoto:  ad.HasPhoto,
	}
	if ad.ID > 0 {
		out.Href = "/edit-ad/" + strconv.FormatInt(ad.ID, 10)
	}
	if ad.Neighborhood != "" {
		if out.Location != "" {
			out.Location += "، " + ad.Neighborhood
		} else {
			out.Location = ad.Neighborhood
		}
	}
	if ad.PriceAmount != nil {
		out.Price = formatAmount(*ad.PriceAmount) + " " + t.Search.Currency
	} else {
		out.Price = t.Search.Negotiable
	}
	if ad.PublishedAt != nil {
		if ts, err := time.Parse(time.RFC3339, *ad.PublishedAt); err == nil {
			out.PublishedAt = ts.Format("2006-01-02")
		}
	}
	return out
}

func formatAmount(v int64) string {
	s := strconv.FormatInt(v, 10)
	neg := strings.HasPrefix(s, "-")
	if neg {
		s = s[1:]
	}
	var b strings.Builder
	for i, r := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}

func (h *PageHandler) MarkedAds(c *gin.Context) {
	h.renderProtected(c, "myinfo_marked_ads", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.MarkedAds)
}

func (h *PageHandler) Setting(c *gin.Context) {
	t := h.i18n.MessagesFor(i18n.FromContext(c))
	pageData := i18n.BuildPage(h.i18n, h.cities, i18n.FromContext(c), h.config.AppName, i18n.CityFromContext(c), c.Request.URL.Path, i18n.LocationsFromContext(c))
	pageData.DefaultCountryCode = h.config.DefaultCountryCode
	pageData.Title = h.config.AppName + " — " + t.Appearance.Title
	pageData.Heading = t.Appearance.Title
	if user, ok := h.currentUser(c); ok {
		pageData.IsAuthenticated = true
		pageData.SessionUserName = user.Name
		pageData.SessionUserMobile = user.Mobile
	}
	c.HTML(http.StatusOK, "myinfo_setting", viewmodel.SettingPage{
		Page: pageData,
		Seas: viewmodel.SeasFor(t),
	})
}

type meResponse struct {
	Authenticated bool         `json:"authenticated"`
	User          *sessionUser `json:"user"`
}

type sessionUser struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Mobile     string `json:"mobile"`
	NationalId string `json:"national_id"`
}

type userAdsResponse struct {
	Ads []userAdJSON `json:"ads"`
}

type userAdJSON struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	PriceAmount  *int64  `json:"price_amount"`
	PriceType    string  `json:"price_type"`
	Currency     string  `json:"currency"`
	CityName     string  `json:"city_name"`
	Neighborhood string  `json:"neighborhood"`
	Thumbnail    string  `json:"thumbnail"`
	HasPhoto     bool    `json:"has_photo"`
	PublishedAt  *string `json:"published_at"`
}
