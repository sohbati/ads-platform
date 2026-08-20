package handler

import (
	"net/http"

	"ads-platform-ui/internal/business/category/service"
	"ads-platform-ui/internal/business/category/viewmodel"
	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/cities"
	"ads-platform-ui/internal/core/config"
	"ads-platform-ui/internal/core/i18n"
	"ads-platform-ui/internal/web/page"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	config  *config.Config
	i18n    *i18n.Registry
	cities  *cities.Catalog
	service service.CategoryService
}

func NewPageHandler(cfg *config.Config, reg *i18n.Registry, catalog *cities.Catalog, svc service.CategoryService) *PageHandler {
	return &PageHandler{config: cfg, i18n: reg, cities: catalog, service: svc}
}

func (h *PageHandler) Index(c *gin.Context) {
	heading := h.i18n.MessagesFor(i18n.FromContext(c)).Nav.Category
	base := page.Base(c, h.config, h.i18n, h.cities, h.config.AppName+" — "+heading, heading)

	vm := viewmodel.CategoryPage{Page: base}

	items, err := h.service.List(c.Request.Context())
	if err != nil {
		vm.LoadError = err.Error()
		c.HTML(http.StatusOK, "category", vm)
		return
	}

	vm.Groups = buildGroups(items)
	if len(vm.Groups) > 0 {
		vm.ActiveID = vm.Groups[0].ID
	}
	c.HTML(http.StatusOK, "category", vm)
}

func buildGroups(items []cdn.Category) []viewmodel.CategoryGroup {
	childrenByParent := make(map[int][]cdn.Category)
	var roots []cdn.Category
	for _, item := range items {
		if item.Parent == nil {
			roots = append(roots, item)
			continue
		}
		childrenByParent[*item.Parent] = append(childrenByParent[*item.Parent], item)
	}

	groups := make([]viewmodel.CategoryGroup, 0, len(roots))
	for _, root := range roots {
		group := viewmodel.CategoryGroup{
			CategoryItem: toItem(root),
		}
		for _, child := range childrenByParent[root.ID] {
			group.Children = append(group.Children, toItem(child))
		}
		group.Columns = chunkColumns(group.Title, group.Href, group.Children, 7)
		groups = append(groups, group)
	}
	return groups
}

// chunkColumns spreads links across vertical columns.
func chunkColumns(title, href string, links []viewmodel.CategoryItem, perCol int) []viewmodel.CategoryColumn {
	if len(links) == 0 {
		return nil
	}
	if perCol < 1 {
		perCol = 7
	}
	cols := make([]viewmodel.CategoryColumn, 0, (len(links)+perCol-1)/perCol)
	for i := 0; i < len(links); i += perCol {
		end := i + perCol
		if end > len(links) {
			end = len(links)
		}
		colTitle := title
		colHref := href
		if i > 0 {
			colTitle = ""
			colHref = ""
		}
		cols = append(cols, viewmodel.CategoryColumn{
			Title: colTitle,
			Href:  colHref,
			Links: links[i:end],
		})
	}
	return cols
}

func toItem(item cdn.Category) viewmodel.CategoryItem {
	return viewmodel.CategoryItem{
		ID:    item.ID,
		Title: item.Title,
		Slug:  item.Slug,
		Href:  "/query-ads?category=" + item.Slug,
		Icon:  iconForSlug(item.Slug),
	}
}

func iconForSlug(slug string) string {
	switch slug {
	case "real-estate", "residential-sale", "residential-rent":
		return "home"
	case "vehicles", "cars", "motorcycles":
		return "car"
	case "digital", "mobile-tablet", "laptop":
		return "device"
	case "home", "appliances":
		return "sofa"
	case "services", "education":
		return "wrench"
	case "jobs", "admin-jobs":
		return "briefcase"
	case "personal", "fashion":
		return "shirt"
	case "leisure":
		return "ball"
	case "social":
		return "users"
	case "industrial":
		return "factory"
	default:
		return "grid"
	}
}
