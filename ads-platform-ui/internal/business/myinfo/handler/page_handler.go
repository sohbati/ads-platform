package handler

import (
	"encoding/json"
	"net/http"

	"ads-platform-ui/internal/core/bff"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"

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

func (h *PageHandler) renderProtected(c *gin.Context, templateName, heading string) {
	user, ok := h.currentUser(c)
	if !ok {
		c.Redirect(http.StatusFound, "/login?next="+c.Request.URL.Path)
		return
	}

	pageData := i18n.BuildPage(h.i18n, h.cities, i18n.FromContext(c), h.config.AppName, i18n.CityFromContext(c), c.Request.URL.Path)
	pageData.Title = h.config.AppName + " — " + heading
	pageData.Heading = heading
	pageData.IsAuthenticated = true
	pageData.SessionUserName = user.Name
	pageData.SessionUserMobile = user.Mobile
	c.HTML(http.StatusOK, templateName, pageData)
}

func (h *PageHandler) Index(c *gin.Context) {
	h.renderProtected(c, "myinfo", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.MyInfo)
}

func (h *PageHandler) UserDetails(c *gin.Context) {
	h.renderProtected(c, "myinfo_user_details", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.UserDetails)
}

func (h *PageHandler) UserAds(c *gin.Context) {
	h.renderProtected(c, "myinfo_user_ads", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.UserAds)
}

func (h *PageHandler) MarkedAds(c *gin.Context) {
	h.renderProtected(c, "myinfo_marked_ads", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.MarkedAds)
}

func (h *PageHandler) Setting(c *gin.Context) {
	h.renderProtected(c, "myinfo_setting", h.i18n.MessagesFor(i18n.FromContext(c)).Nav.Setting)
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
