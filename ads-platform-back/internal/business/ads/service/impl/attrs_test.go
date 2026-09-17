package impl

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	searchclient "ads-platform/internal/business/search/client"
)

func residentialSaleSchemaJSON() json.RawMessage {
	return json.RawMessage(`{
		"type":"object",
		"additionalProperties":false,
		"properties":{
			"area_m2":{"type":"integer","minimum":1,"maximum":100000},
			"rooms":{"type":"integer","minimum":0,"maximum":50},
			"has_storage":{"type":"boolean"},
			"storage_m2":{"type":"number","minimum":0,"maximum":10000},
			"floors":{"type":"integer","minimum":1,"maximum":200},
			"units":{"type":"integer","minimum":1,"maximum":10000},
			"orientation":{"type":"string","enum":["north","south"]},
			"balcony":{"type":"boolean"},
			"flooring":{"type":"string","enum":["ceramic","parquet","laminate","mosaic"]},
			"carpeted":{"type":"boolean"},
			"toilets_iranian":{"type":"integer","enum":[0,1,2,3,4,5,6],"minimum":0,"maximum":6},
			"toilets_western":{"type":"integer","enum":[0,1,2,3,4,5,6],"minimum":0,"maximum":6},
			"bathrooms":{"type":"integer","enum":[0,1,2,3,4,5,6],"minimum":0,"maximum":6},
			"hot_water":{"type":"string","enum":["water_heater","package","boiler_room"]},
			"facade":{"type":"string","enum":["stone","brick","cement","glass","other"]},
			"deed_type":{"type":"string","enum":["single","joint","cooperative","booklet","promissory"]}
		},
		"required":["area_m2","rooms"]
	}`)
}

func residentialSaleCatalog() *fakeCatalog {
	name := "residential-sale"
	return &fakeCatalog{
		categories: []searchclient.Category{{
			ID:                             11,
			Slug:                           "residential-sale",
			Title:                          "Residential sale",
			AdsAttrsJSONSchemaTemplateName: &name,
			DescendantIDs:                  []int{11},
		}},
		cities: []searchclient.City{{ID: 1, Slug: "tehran", Name: "Tehran"}},
		schemas: []searchclient.AttrSchema{{
			Name:       "residential-sale",
			Title:      "Residential sale",
			JSONSchema: residentialSaleSchemaJSON(),
		}},
	}
}

func TestFilterAttrsBySchemaKeepsKnownResidentialSaleFields(t *testing.T) {
	raw := json.RawMessage(`{
		"area_m2":90,
		"rooms":2,
		"has_storage":true,
		"storage_m2":8.5,
		"floors":5,
		"units":10,
		"orientation":"north",
		"balcony":true,
		"flooring":"ceramic",
		"carpeted":false,
		"toilets_iranian":1,
		"toilets_western":1,
		"bathrooms":2,
		"hot_water":"package",
		"facade":"stone",
		"deed_type":"booklet",
		"unknown":"drop-me"
	}`)
	got, err := filterAttrsBySchema(raw, searchclient.AttrSchema{JSONSchema: residentialSaleSchemaJSON()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var attrs map[string]any
	if err := json.Unmarshal(got, &attrs); err != nil {
		t.Fatal(err)
	}
	if _, ok := attrs["unknown"]; ok {
		t.Fatalf("unknown key stored: %s", got)
	}
	if attrs["area_m2"] != float64(90) || attrs["rooms"] != float64(2) {
		t.Fatalf("required ints: %s", got)
	}
	if attrs["has_storage"] != true || attrs["carpeted"] != false || attrs["storage_m2"] != 8.5 {
		t.Fatalf("bools/number: %s", got)
	}
	if attrs["orientation"] != "north" || attrs["deed_type"] != "booklet" || attrs["hot_water"] != "package" {
		t.Fatalf("enums: %s", got)
	}
}

func TestFilterAttrsBySchemaRejectsInvalidAndMissingRequired(t *testing.T) {
	schema := searchclient.AttrSchema{JSONSchema: residentialSaleSchemaJSON()}
	_, err := filterAttrsBySchema(json.RawMessage(`{"rooms":2}`), schema)
	if err == nil {
		t.Fatal("expected missing area_m2")
	}
	_, err = filterAttrsBySchema(json.RawMessage(`{"area_m2":90,"rooms":2,"orientation":"west"}`), schema)
	if err == nil {
		t.Fatal("expected invalid enum")
	}
	_, err = filterAttrsBySchema(json.RawMessage(`{"area_m2":90,"rooms":2,"floors":0}`), schema)
	if err == nil {
		t.Fatal("expected floors below min")
	}
	_, err = filterAttrsBySchema(json.RawMessage(`{"area_m2":90,"rooms":2,"bathrooms":7}`), schema)
	if err == nil {
		t.Fatal("expected bathrooms above max")
	}
}

func TestCreateStoresResidentialSaleAttrs(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, residentialSaleCatalog(), nil, 8, 10<<20)
	in := validInput()
	in.CategoryID = 11
	in.Attrs = json.RawMessage(`{
		"area_m2":120,
		"rooms":3,
		"has_storage":true,
		"storage_m2":6,
		"floors":4,
		"units":8,
		"orientation":"south",
		"balcony":false,
		"flooring":"parquet",
		"carpeted":true,
		"toilets_iranian":1,
		"toilets_western":0,
		"bathrooms":1,
		"hot_water":"boiler_room",
		"facade":"brick",
		"deed_type":"promissory",
		"property_price":999
	}`)
	ad, err := svc.Create(context.Background(), in)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	body := string(ad.Attrs)
	for _, want := range []string{
		`"area_m2":120`, `"rooms":3`, `"has_storage":true`, `"storage_m2":6`,
		`"orientation":"south"`, `"flooring":"parquet"`, `"hot_water":"boiler_room"`,
		`"facade":"brick"`, `"deed_type":"promissory"`, `"carpeted":true`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("attrs missing %s: %s", want, body)
		}
	}
	if strings.Contains(body, "property_price") {
		t.Fatalf("dropped field stored: %s", body)
	}

	got, err := svc.GetPublic(context.Background(), ad.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.CategoryID != 11 {
		t.Fatalf("public category_id=%d", got.CategoryID)
	}
	if !strings.Contains(string(got.Attrs), `"orientation":"south"`) {
		t.Fatalf("public attrs=%s", got.Attrs)
	}
}

func TestCreateRejectsInvalidResidentialSaleAttrs(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, residentialSaleCatalog(), nil, 8, 10<<20)
	in := validInput()
	in.CategoryID = 11
	in.Attrs = json.RawMessage(`{"area_m2":90}`)
	_, err := svc.Create(context.Background(), in)
	assertAppErrorCode(t, err, "AD_INVALID_ATTRS")
}

func TestUpdateStoresResidentialSaleAttrs(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, residentialSaleCatalog(), nil, 8, 10<<20)
	in := validInput()
	in.CategoryID = 11
	in.Attrs = json.RawMessage(`{"area_m2":80,"rooms":1}`)
	ad, err := svc.Create(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	in.Attrs = json.RawMessage(`{"area_m2":95,"rooms":2,"balcony":true,"facade":"glass"}`)
	updated, err := svc.Update(context.Background(), ad.ID, in)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	body := string(updated.Attrs)
	if !strings.Contains(body, `"balcony":true`) || !strings.Contains(body, `"facade":"glass"`) {
		t.Fatalf("updated attrs=%s", body)
	}
}

func TestCreateWithoutSchemaTemplateKeepsPostedAttrs(t *testing.T) {
	svc := NewAdService(newFakeAdRepo(), &fakeImageRepo{}, leafCatalog(), nil, 8, 10<<20)
	ad, err := svc.Create(context.Background(), validInput())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(ad.Attrs), `"rooms":2`) {
		t.Fatalf("attrs=%s", ad.Attrs)
	}
}
