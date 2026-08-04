package router

import (
	"html/template"
	"net/http"

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

	tmpl := template.Must(template.New("").Funcs(coreview.FuncMap()).ParseGlob("templates/**/*.gohtml"))
	router.SetHTMLTemplate(tmpl)

	router.Static("/static", "./static")

	router.NoRoute(func(c *gin.Context) {
		cfg := r.container.Config
		loc := i18n.FromContext(c)
		city := i18n.CityFromContext(c)
		page := i18n.BuildPage(r.container.I18n, r.container.Cities, loc, cfg.AppName, city, c.Request.URL.Path)
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
	}

	return router
}
