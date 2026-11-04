package impl

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"ads-platform/internal/business/adreview/model"
)

// minTitleSearch is the shortest title search; the trigram index cannot help with fewer letters.
const minTitleSearch = 3

var digitReplacer = strings.NewReplacer(
	"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4", "۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
	"٠", "0", "١", "1", "٢", "2", "٣", "3", "٤", "4", "٥", "5", "٦", "6", "٧", "7", "٨", "8", "٩", "9",
)

// An Iranian mobile in any common form; the group is the 10 digits after the country code.
var mobilePattern = regexp.MustCompile(`^(?:\+98|0098|98|0)?(9\d{9})$`)

var allDigits = regexp.MustCompile(`^\d+$`)

// searchTerms is what one search box can mean: exactly one field is set.
type searchTerms struct {
	adID    int64
	mobiles []string
	title   string
}

// parseSearch reads a search box: a mobile number, an ad ID (digits, optionally after #), or
// part of a title. Persian and Arabic digits count as digits.
func parseSearch(raw string) (searchTerms, bool) {
	text := strings.TrimSpace(digitReplacer.Replace(raw))
	if text == "" {
		return searchTerms{}, true
	}
	compact := strings.NewReplacer(" ", "", "-", "").Replace(text)
	if match := mobilePattern.FindStringSubmatch(compact); match != nil {
		// Customers are stored with either prefix, so match every form.
		national := match[1]
		return searchTerms{mobiles: []string{"+98" + national, "0" + national, "98" + national}}, true
	}
	if id := strings.TrimPrefix(compact, "#"); allDigits.MatchString(id) {
		parsed, err := strconv.ParseInt(id, 10, 64)
		if err != nil || parsed <= 0 {
			return searchTerms{}, false
		}
		return searchTerms{adID: parsed}, true
	}
	if utf8.RuneCountInString(text) < minTitleSearch {
		return searchTerms{}, false
	}
	return searchTerms{title: text}, true
}

func encodeCursor(cursor model.ReviewCursor) string {
	raw, _ := json.Marshal(cursor)
	return base64.RawURLEncoding.EncodeToString(raw)
}

func decodeCursor(value string) (*model.ReviewCursor, bool) {
	if value == "" {
		return nil, true
	}
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, false
	}
	var cursor model.ReviewCursor
	if err := json.Unmarshal(raw, &cursor); err != nil || cursor.ID <= 0 || cursor.At.IsZero() {
		return nil, false
	}
	cursor.At = cursor.At.UTC()
	return &cursor, true
}

// cursorAfter is the sort position of the given row, matching the order List uses.
func cursorAfter(item model.ReviewListItem, pending bool) model.ReviewCursor {
	at := item.UpdatedAt
	if pending && item.PendingSince != nil {
		at = *item.PendingSince
	}
	return model.ReviewCursor{At: at.UTC().Truncate(time.Microsecond), ID: item.ID}
}
