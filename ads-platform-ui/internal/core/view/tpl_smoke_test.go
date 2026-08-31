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
		Page: i18n.Page{
			Title:           "t",
			CityDisplayName: "Tehran",
			T: i18n.Messages{
				Hero: i18n.HeroMessages{Title: "On the surface, in %s", Subtitle: "sub"},
			},
		},
		Search: &viewmodel.SearchResults{},
	}
	search := viewmodel.QueryAdsPage{
		Page: i18n.Page{
			Title:       "t",
			SearchQuery: "bike",
			T: i18n.Messages{
				Hero: i18n.HeroMessages{Title: "On the surface, in %s", Subtitle: "sub"},
			},
		},
		Search: &viewmodel.SearchResults{
			Query:      "bike",
			Total:      2,
			Page:       1,
			TotalPages: 2,
			NextURL:    "/query-ads?q=bike&page=2",
			Ads: []viewmodel.SearchAd{
				{Title: "Bike", Price: "1,000 IRR", Location: "Tehran", Thumbnail: "/t.jpg", HasPhoto: true, PublishedAt: "2026-08-01", Href: "/ad/1"},
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
		if !bytes.Contains(buf.Bytes(), []byte("brand__wordmark")) || !bytes.Contains(buf.Bytes(), []byte(`class="brand__tld">.ir`)) {
			t.Fatalf("query_ads (%s): expected ruab.ir brand lockup", name)
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
					Intro:         "intro",
					PicturesHint:  "Up to %d photos",
					PicturesAdd:   "Add photos",
					PictureRemove: "Remove photo",
					PictureView:   "View photo",
					PictureClose:  "Close",
					Submit:        "Publish",
				},
				AdDetail: i18n.AdDetailMessages{
					PrevPhoto:    "Previous photo",
					NextPhoto:    "Next photo",
					PhotoCounter: "%d / %d",
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

func TestMyInfoSettingTemplateRenders(t *testing.T) {
	tmpl, err := template.New("").Funcs(FuncMap(nil)).ParseGlob("../../../templates/**/*.gohtml")
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}

	page := myinfovm.SettingPage{
		Page: i18n.Page{
			Title:           "t",
			Heading:         "Appearance",
			IsAuthenticated: true,
			T: i18n.Messages{
				Appearance: i18n.AppearanceMessages{
					Title:       "Appearance",
					Description: "Light or dark",
					GroupAria:   "theme",
					Light:       "Light",
					Dark:        "Dark",
				},
				Nav: i18n.NavMessages{Setting: "Settings", Logout: "Out"},
			},
		},
		Themes: []myinfovm.ThemeOption{
			{ID: "light", Name: "Light"},
			{ID: "tide", Name: "Tide"},
			{ID: "dark", Name: "Dark"},
		},
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "myinfo_setting", page); err != nil {
		t.Fatalf("execute myinfo_setting: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("data-theme-id=\"tide\"")) {
		t.Fatalf("expected tide swatch markup, got %s", buf.String())
	}
}

func TestAdDetailTemplateRenders(t *testing.T) {
	tmpl, err := template.New("").Funcs(FuncMap(nil)).ParseGlob("../../../templates/**/*.gohtml")
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}

	withPhotos := viewmodel.AdDetailPage{
		Page: i18n.Page{
			Title: "t",
			T: i18n.Messages{
				AdDetail: i18n.AdDetailMessages{
					PrevPhoto:    "Previous photo",
					NextPhoto:    "Next photo",
					PhotoCounter: "%d / %d",
				},
			},
		},
		Ad: &viewmodel.AdDetail{
			Title:       "Bike",
			Price:       "1,000 IRR",
			Location:    "Tehran",
			PublishedAt: "2026-08-01",
			Description: "Nice bike",
			Images:      []string{"/a.webp", "/b.webp"},
			HasPhone:    true,
			PhoneMasked: "09*********",
		},
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "ad_detail", withPhotos); err != nil {
		t.Fatalf("execute ad_detail: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("/a.webp")) || !bytes.Contains(buf.Bytes(), []byte("data-ad-next")) {
		t.Fatalf("expected gallery markup, got %s", buf.String())
	}
	if !bytes.Contains(buf.Bytes(), []byte("09*********")) || !bytes.Contains(buf.Bytes(), []byte("data-ad-show-phone")) {
		t.Fatalf("expected masked phone markup, got %s", buf.String())
	}
	if bytes.Contains(buf.Bytes(), []byte("09121110001")) {
		t.Fatalf("full phone must not be in HTML: %s", buf.String())
	}

	notFound := viewmodel.AdDetailPage{
		Page: i18n.Page{
			Title: "t",
			T:     i18n.Messages{AdDetail: i18n.AdDetailMessages{NotFound: "missing"}},
		},
		NotFound: true,
	}
	buf.Reset()
	if err := tmpl.ExecuteTemplate(&buf, "ad_detail", notFound); err != nil {
		t.Fatalf("execute ad_detail not found: %v", err)
	}
}
