package media

import "testing"

func TestPublicURL(t *testing.T) {
	cases := []struct {
		base, path, want string
	}{
		{"http://localhost:8098", "/ads-media/ads/1/1_1.webp", "http://localhost:8098/ads-media/ads/1/1_1.webp"},
		{"http://localhost:8098/", "ads-media/ads/1/1_1.webp", "http://localhost:8098/ads-media/ads/1/1_1.webp"},
		{"http://localhost:8098", "https://cdn.example/x.webp", "https://cdn.example/x.webp"},
		{"", "/ads-media/ads/1/1_1.webp", "/ads-media/ads/1/1_1.webp"},
		{"http://localhost:8098", "", ""},
	}
	for _, tc := range cases {
		if got := PublicURL(tc.base, tc.path); got != tc.want {
			t.Errorf("PublicURL(%q, %q) = %q, want %q", tc.base, tc.path, got, tc.want)
		}
	}
}
