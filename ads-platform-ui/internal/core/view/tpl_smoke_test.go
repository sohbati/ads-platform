package view

import (
	"bytes"
	"html/template"
	"testing"

	"ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/i18n"
)

func TestQueryAdsTemplateRenders(t *testing.T) {
	tmpl, err := template.New("").Funcs(FuncMap()).ParseGlob("../../../templates/**/*.gohtml")
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}

	landing := viewmodel.QueryAdsPage{Page: i18n.Page{Title: "t"}}
	search := viewmodel.QueryAdsPage{
		Page: i18n.Page{Title: "t", SearchQuery: "bike"},
		Search: &viewmodel.SearchResults{
			Query:      "bike",
			Total:      2,
			Page:       1,
			TotalPages: 2,
			NextURL:    "/query-ads?q=bike&page=2",
			Ads: []viewmodel.SearchAd{
				{Title: "Bike", Price: "1,000 IRR", Location: "Tehran", Thumbnail: "/t.jpg", HasPhoto: true, PublishedAt: "2026-08-01"},
				{Title: "Old bike", Price: "Negotiable"},
			},
		},
	}

	for name, data := range map[string]viewmodel.QueryAdsPage{"landing": landing, "search": search} {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "query_ads", data); err != nil {
			t.Fatalf("execute query_ads (%s): %v", name, err)
		}
	}
}
