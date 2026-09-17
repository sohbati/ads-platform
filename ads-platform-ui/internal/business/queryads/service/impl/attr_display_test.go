package impl

import (
	"encoding/json"
	"testing"

	"ads-platform-ui/internal/core/cdn"
	"ads-platform-ui/internal/core/i18n"
)

func TestLabeledAttrsResidentialSale(t *testing.T) {
	name := "residential-sale"
	categories := []cdn.Category{{
		ID:                             11,
		AdsAttrsJSONSchemaTemplateName: &name,
	}}
	schemas := []cdn.AttrSchema{{
		Name: "residential-sale",
		JSONSchema: json.RawMessage(`{
			"properties":{
				"area_m2":{"type":"integer","title":"متراژ (متر مربع)"},
				"has_storage":{"type":"boolean","title":"انباری دارد؟"},
				"orientation":{"type":"string","title":"شمالی / جنوبی؟","x-enumVocab":"orientation"},
				"storage_m2":{"type":"number","title":"متراژ انباری"}
			}
		}`),
	}}
	enums := json.RawMessage(`{"orientation":{"north":{"fa":"شمالی","en":"North"}}}`)
	attrs := json.RawMessage(`{"area_m2":90,"has_storage":true,"orientation":"north","storage_m2":8.5,"skip":1}`)
	tmsg := i18n.Messages{AdDetail: i18n.AdDetailMessages{Yes: "بله", No: "خیر"}}

	got := labeledAttrs(i18n.FA, tmsg, 11, attrs, categories, schemas, enums)
	if len(got) != 4 {
		t.Fatalf("got %d rows: %+v", len(got), got)
	}
	if got[0].Label != "متراژ (متر مربع)" || got[0].Value != "90" {
		t.Fatalf("area: %+v", got[0])
	}
	if got[1].Value != "بله" {
		t.Fatalf("storage: %+v", got[1])
	}
	if got[2].Value != "شمالی" {
		t.Fatalf("orientation: %+v", got[2])
	}
	if got[3].Value != "8.5" {
		t.Fatalf("storage_m2: %+v", got[3])
	}
}
