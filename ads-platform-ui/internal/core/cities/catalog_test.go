package cities

import "testing"

func TestExpandToCityIDsAndSearchPlace(t *testing.T) {
	cat, err := parseRecords([]byte(`[
		{"id":1,"name":"Iran","slug":"iran","type":"0"},
		{"id":10,"name":"Tehran","slug":"tehran-province","parent":1,"type":"1"},
		{"id":11,"name":"Tehran","slug":"tehran","parent":10,"type":"2"},
		{"id":12,"name":"Pardis","slug":"pardis","parent":10,"type":"2"},
		{"id":20,"name":"Karaj","slug":"karaj","parent":1,"type":"2"}
	]`), "test")
	if err != nil {
		t.Fatal(err)
	}

	ids := cat.ExpandToCityIDs([]string{"tehran-province"})
	if len(ids) != 2 || ids[0] != 11 || ids[1] != 12 {
		t.Fatalf("province expand = %v, want [11 12]", ids)
	}

	place, csv := cat.SearchPlace([]string{"tehran-province"}, "tehran")
	if place != "iran" || csv != "11,12" {
		t.Fatalf("SearchPlace(province) = %s %s, want iran 11,12", place, csv)
	}

	place, csv = cat.SearchPlace([]string{"tehran"}, "tehran")
	if place != "tehran" || csv != "" {
		t.Fatalf("SearchPlace(city) = %s %q, want tehran \"\"", place, csv)
	}

	if got := cat.PrimaryCitySlug([]string{"tehran-province"}, "karaj"); got != "tehran" {
		t.Fatalf("PrimaryCitySlug = %s, want tehran", got)
	}

	slugs := cat.ParseLocationSlugs("tehran, pardis, unknown, tehran")
	if len(slugs) != 2 || slugs[0] != "tehran" || slugs[1] != "pardis" {
		t.Fatalf("ParseLocationSlugs = %v", slugs)
	}
}
