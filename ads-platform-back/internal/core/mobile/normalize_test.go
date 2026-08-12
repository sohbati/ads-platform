package mobile

import "testing"

func TestNormalizeIranMobile(t *testing.T) {
	n := NewNormalizer("+98")

	cases := []struct {
		in   string
		want string
	}{
		{"09125697463", "+989125697463"},
		{"9125697463", "+989125697463"},
		{"989125697463", "+989125697463"},
		{"+989125697463", "+989125697463"},
		{"00989125697463", "+989125697463"},
		{"0098 912 569 7463", "+989125697463"},
		{"+98 912 569 7463", "+989125697463"},
		{"۰۹۱۲۵۶۹۷۴۶۳", "+989125697463"},
	}

	for _, tc := range cases {
		got, err := n.Normalize(tc.in)
		if err != nil {
			t.Fatalf("Normalize(%q): unexpected error: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("Normalize(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestNormalizeEmpty(t *testing.T) {
	n := NewNormalizer("+98")
	if _, err := n.Normalize("   "); err == nil {
		t.Fatal("expected error for empty mobile")
	}
}

func TestNormalizeInvalid(t *testing.T) {
	n := NewNormalizer("+98")
	if _, err := n.Normalize("123"); err == nil {
		t.Fatal("expected error for invalid mobile")
	}
}
