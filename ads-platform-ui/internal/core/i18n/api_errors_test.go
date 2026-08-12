package i18n

import "testing"

func TestResolveAPIError(t *testing.T) {
	catalog := map[string]string{
		"_default":          "Something went wrong.",
		"OTP_VERIFY_FAILED": "The code is incorrect.",
		"OTP_EXPIRED":       "Code expired for %s.",
	}

	if got := ResolveAPIError(catalog, "OTP_VERIFY_FAILED", nil, ""); got != "The code is incorrect." {
		t.Fatalf("unexpected message: %q", got)
	}

	if got := ResolveAPIError(catalog, "OTP_EXPIRED", []string{"0912"}, ""); got != "Code expired for 0912." {
		t.Fatalf("unexpected message: %q", got)
	}

	if got := ResolveAPIError(catalog, "UNKNOWN", nil, ""); got != "Something went wrong." {
		t.Fatalf("expected default, got %q", got)
	}
}
