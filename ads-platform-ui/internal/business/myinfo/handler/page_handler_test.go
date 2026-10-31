package handler

import (
	"testing"

	"ads-platform-ui/internal/core/i18n"
)

func TestCustomerAdStatus(t *testing.T) {
	registry, err := i18n.LoadRegistry(i18n.FA)
	if err != nil {
		t.Fatal(err)
	}
	msgs := registry.MessagesFor(i18n.FA)

	cases := []struct {
		ad         customerAdJSON
		wantLabel  string
		wantTone   string
		wantReason string
	}{
		{customerAdJSON{Status: "pending"}, msgs.MyAds.StatusPending, "pending", ""},
		{customerAdJSON{Status: "active"}, msgs.MyAds.StatusActive, "active", ""},
		{customerAdJSON{Status: "draft", RejectionReason: "عکس ندارد"}, msgs.MyAds.StatusRejected, "rejected", "عکس ندارد"},
		{customerAdJSON{Status: "draft"}, msgs.MyAds.StatusDraft, "draft", ""},
		// A resubmitted ad keeps its old reason in the backend; it must not show while pending.
		{customerAdJSON{Status: "pending", RejectionReason: "old"}, msgs.MyAds.StatusPending, "pending", ""},
	}
	for _, c := range cases {
		got := toSearchAd(c.ad, msgs, "")
		if got.StatusLabel != c.wantLabel || got.StatusTone != c.wantTone || got.RejectionReason != c.wantReason {
			t.Errorf("%+v: got label %q tone %q reason %q", c.ad, got.StatusLabel, got.StatusTone, got.RejectionReason)
		}
		if c.wantLabel == "" {
			t.Errorf("%s: empty label in bundle", c.ad.Status)
		}
	}
}
