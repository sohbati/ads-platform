package view

import (
	"bytes"
	"html/template"
	"testing"

	myinfovm "ads-platform-ui/internal/business/myinfo/viewmodel"
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

	edit := newadvm.NewAdPage{
		Page: i18n.Page{
			Title:   "t",
			Heading: "Edit ad",
			T: i18n.Messages{
				NewAd: i18n.NewAdMessages{
					EditIntro:           "edit intro",
					SubmitEdit:          "Save",
					PicturesReplaceHint: "replace photos",
				},
			},
		},
		Bootstrap: newadvm.Bootstrap{Mode: "edit", AdID: 9, MaxPictures: 8, CityName: "Tehran", Enums: []byte(`{}`)},
	}
	buf.Reset()
	if err := tmpl.ExecuteTemplate(&buf, "new_ad", edit); err != nil {
		t.Fatalf("execute new_ad edit: %v", err)
	}
}

func TestMyInfoUserAdsTemplateRenders(t *testing.T) {
	tmpl, err := template.New("").Funcs(FuncMap(nil)).ParseGlob("../../../templates/**/*.gohtml")
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}

	empty := myinfovm.UserAdsPage{
		Page: i18n.Page{
			Title:           "t",
			Heading:         "My ads",
			IsAuthenticated: true,
			T: i18n.Messages{
				MyAds: i18n.MyAdsMessages{Empty: "no ads", PostCta: "post"},
				Nav:   i18n.NavMessages{UserAds: "My ads", Logout: "Out"},
			},
		},
	}
	withAds := myinfovm.UserAdsPage{
		Page: i18n.Page{
			Title:           "t",
			Heading:         "My ads",
			IsAuthenticated: true,
			T: i18n.Messages{
				Nav: i18n.NavMessages{UserAds: "My ads", Logout: "Out"},
			},
		},
		Ads: []viewmodel.SearchAd{
			{ID: 9, Href: "/edit-ad/9", Title: "Bike", Price: "1,000 IRR", Location: "Tehran", Thumbnail: "/t.jpg", PublishedAt: "2026-08-01"},
		},
	}
	unavailable := myinfovm.UserAdsPage{
		Page: i18n.Page{
			Title:           "t",
			Heading:         "My ads",
			IsAuthenticated: true,
			T: i18n.Messages{
				MyAds: i18n.MyAdsMessages{Unavailable: "down"},
				Nav:   i18n.NavMessages{UserAds: "My ads", Logout: "Out"},
			},
		},
		Unavailable: true,
	}

	for name, data := range map[string]any{
		"empty":       empty,
		"ads":         withAds,
		"unavailable": unavailable,
	} {
		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "myinfo_user_ads", data); err != nil {
			t.Fatalf("execute myinfo_user_ads (%s): %v", name, err)
		}
	}
}
