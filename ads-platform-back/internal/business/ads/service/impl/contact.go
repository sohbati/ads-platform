package impl

import (
	"encoding/json"
	"strings"
)

type contactPayload struct {
	Phone     string `json:"phone"`
	ShowPhone *bool  `json:"show_phone"`
}

func contactPhone(raw json.RawMessage) (string, bool) {
	if len(raw) == 0 {
		return "", false
	}
	var c contactPayload
	if json.Unmarshal(raw, &c) != nil {
		return "", false
	}
	if c.ShowPhone != nil && !*c.ShowPhone {
		return "", false
	}
	phone := strings.TrimSpace(c.Phone)
	if phone == "" {
		return "", false
	}
	return phone, true
}

func maskPhone(phone string) string {
	digits := nationalDigits(phone)
	if digits == "" {
		return "***"
	}
	if len(digits) <= 2 {
		return strings.Repeat("*", len(digits))
	}
	return digits[:2] + strings.Repeat("*", len(digits)-2)
}

func displayPhone(phone string) string {
	if digits := nationalDigits(phone); digits != "" {
		return digits
	}
	return strings.TrimSpace(phone)
}

func nationalDigits(phone string) string {
	var b strings.Builder
	for _, r := range phone {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= '۰' && r <= '۹':
			b.WriteRune('0' + (r - '۰'))
		case r >= '٠' && r <= '٩':
			b.WriteRune('0' + (r - '٠'))
		}
	}
	digits := b.String()
	switch {
	case strings.HasPrefix(digits, "98") && len(digits) >= 12:
		return "0" + digits[2:]
	case len(digits) == 10 && digits[0] != '0':
		return "0" + digits
	default:
		return digits
	}
}
