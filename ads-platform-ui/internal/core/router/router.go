package router

import (
	"html/template"
	"log"
	"net/http"

	"ads-platform-ui/internal/core/assets"
	appContainer "ads-platform-ui/internal/core/container"
	"ads-platform-ui/internal/core/i18n"
	coreview "ads-platform-ui/internal/core/view"

	"github.com/gin-gonic/gin"
)

type Router struct {
	container *appContainer.AppContainer
}

func NewRouter(c *appContainer.AppContainer) *Router {
	return &Router{container: c}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(i18n.Middleware(r.container.I18n, r.container.Cities, r.container.Config.DefaultCity))

	// Content-hashed asset URLs let /static be served as immutable: browsers
	// cache files for a year and never revalidate; changed content gets a new URL.
	var assetURL func(string) string
	if manifest, err := assets.Load("./static", "/static"); err != nil {
		log.Printf("assets: falling back to unversioned URLs: %v", err)
	} else {
		assetURL = manifest.URL
	}

	tmpl := template.Must(template.New("").Funcs(coreview.FuncMap(assetURL)).ParseGlob("templates/**/*.gohtml"))
	router.SetHTMLTemplate(tmpl)

	static := router.Group("/static", func(c *gin.Context) {
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
	})
	static.Static("/", "./static")

	router.NoRoute(func(c *gin.Context) {
		cfg := r.container.Config
		loc := i18n.FromContext(c)
		city := i18n.CityFromContext(c)
		page := i18n.BuildPage(r.container.I18n, r.container.Cities, loc, cfg.AppName, city, c.Request.URL.Path, i18n.LocationsFromContext(c))
		page.DefaultCountryCode = cfg.DefaultCountryCode
		page.Title = cfg.AppName + " — " + page.T.Error.Title
		c.HTML(http.StatusNotFound, "error_404", page)
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ads-platform-ui"})
	})

	// Site map: / → query-ads
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/query-ads")
	})

	router.GET("/query-ads", r.container.QueryAds.PageHandler.Index)

	router.GET("/login", r.container.Auth.PageHandler.Login)

	router.GET("/my-info", r.container.MyInfo.PageHandler.Index)
	router.GET("/my-info/user-details", r.container.MyInfo.PageHandler.UserDetails)
	router.GET("/my-info/user-ads", r.container.MyInfo.PageHandler.UserAds)
	router.GET("/my-info/marked-ads", r.container.MyInfo.PageHandler.MarkedAds)
	router.GET("/my-info/setting", r.container.MyInfo.PageHandler.Setting)

	router.GET("/new-ad", r.container.NewAd.PageHandler.Index)
	router.GET("/category", r.container.Category.PageHandler.Index)
	router.GET("/location", r.container.Location.PageHandler.Index)

	api := router.Group("/api/v1")
	{
		api.GET("/categories", r.container.Category.APIHandler.List)
		api.GET("/cities", r.container.Location.APIHandler.ListCities)
		api.POST("/ads", r.container.NewAd.APIHandler.Create)

		auth := api.Group("/auth")
		{
			auth.POST("/otp/:mobile/send", r.container.Auth.APIHandler.SendOTP)
			auth.POST("/otp/:mobile/verify", r.container.Auth.APIHandler.VerifyOTP)
			auth.GET("/me", r.container.Auth.APIHandler.Me)
			auth.POST("/logout", r.container.Auth.APIHandler.Logout)
		}
	}

	return router
}
