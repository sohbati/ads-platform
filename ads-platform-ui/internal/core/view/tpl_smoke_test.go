package view

import (
	"bytes"
	"html/template"
	"testing"

	newadvm "ads-platform-ui/internal/business/newad/viewmodel"
	"ads-platform-ui/internal/business/queryads/viewmodel"
	"ads-platform-ui/internal/core/i18n"
)

func TestQueryAdsTemplateRenders(t *testing.T) {
	tmpl, err := template.New("").Funcs(FuncMap(nil)).ParseGlob("../../../templates/**/*.gohtml")
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}

	landing := viewmodel.QueryAdsPage{
		Page:   i18n.Page{Title: "t", CityDisplayName: "Tehran"},
		Search: &viewmodel.SearchResults{},
	}
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

	for name, data := range map[string]any{
		"landing": landing,
		"search":  search,
	} {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "query_ads", data); err != nil {
			t.Fatalf("execute query_ads (%s): %v", name, err)
		}
	}
}

func TestNewAdTemplateRenders(t *testing.T) {
	tmpl, err := template.New("").Funcs(FuncMap(nil)).ParseGlob("../../../templates/**/*.gohtml")
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}

	page := newadvm.NewAdPage{
		Page: i18n.Page{
			Title:           "t",
			Heading:         "New ad",
			CityDisplayName: "Tehran",
			T: i18n.Messages{
				NewAd: i18n.NewAdMessages{
					Intro:        "intro",
					PicturesHint: "Up to %d photos",
					Submit:       "Publish",
				},
			},
		},
		Bootstrap: newadvm.Bootstrap{MaxPictures: 8, CityName: "Tehran", Enums: []byte(`{}`)},
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "new_ad", page); err != nil {
		t.Fatalf("execute new_ad: %v", err)
	}
}
