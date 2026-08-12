package mobile

import (
	"fmt"
	"strings"
	"unicode"
)

var digitMap = map[rune]rune{
	'۰': '0', '۱': '1', '۲': '2', '۳': '3', '۴': '4',
	'۵': '5', '۶': '6', '۷': '7', '۸': '8', '۹': '9',
	'٠': '0', '١': '1', '٢': '2', '٣': '3', '٤': '4',
	'٥': '5', '٦': '6', '٧': '7', '٨': '8', '٩': '9',
}

// Normalizer canonicalizes mobile numbers to E.164 using a default country code.
type Normalizer struct {
	DefaultCountryCode string
}

func NewNormalizer(defaultCountryCode string) *Normalizer {
	if defaultCountryCode == "" {
		defaultCountryCode = "+98"
	}
	return &Normalizer{DefaultCountryCode: normalizeCountryCode(defaultCountryCode)}
}

func (n *Normalizer) Normalize(input string) (string, error) {
	defaultCC := n.DefaultCountryCode
	if defaultCC == "" {
		defaultCC = "+98"
	}

	s := strings.TrimSpace(input)
	if s == "" {
		return "", fmt.Errorf("mobile is empty")
	}

	s = convertDigits(s)
	s = stripSeparators(s)

	if strings.HasPrefix(s, "00") {
		s = "+" + s[2:]
	}

	ccDigits := strings.TrimPrefix(defaultCC, "+")

	switch {
	case strings.HasPrefix(s, "+"):
		normalized := "+" + digitsOnly(s[1:])
		if err := validateMobile(normalized, ccDigits); err != nil {
			return "", err
		}
		return normalized, nil
	case strings.HasPrefix(s, "0"):
		normalized := defaultCC + s[1:]
		if err := validateMobile(normalized, ccDigits); err != nil {
			return "", err
		}
		return normalized, nil
	case strings.HasPrefix(s, ccDigits):
		normalized := "+" + s
		if err := validateMobile(normalized, ccDigits); err != nil {
			return "", err
		}
		return normalized, nil
	default:
		normalized := defaultCC + s
		if err := validateMobile(normalized, ccDigits); err != nil {
			return "", err
		}
		return normalized, nil
	}
}

func normalizeCountryCode(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return "+98"
	}
	if !strings.HasPrefix(code, "+") {
		code = "+" + code
	}
	return "+" + digitsOnly(code[1:])
}

func convertDigits(input string) string {
	var b strings.Builder
	b.Grow(len(input))
	for _, r := range input {
		if mapped, ok := digitMap[r]; ok {
			b.WriteRune(mapped)
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func stripSeparators(input string) string {
	var b strings.Builder
	b.Grow(len(input))
	for _, r := range input {
		if unicode.IsSpace(r) || r == '-' || r == '(' || r == ')' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func digitsOnly(input string) string {
	var b strings.Builder
	b.Grow(len(input))
	for _, r := range input {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func validateMobile(normalized, ccDigits string) error {
	if !strings.HasPrefix(normalized, "+") {
		return fmt.Errorf("invalid mobile format")
	}

	digits := digitsOnly(normalized[1:])
	if digits == "" {
		return fmt.Errorf("invalid mobile format")
	}

	if ccDigits == "98" {
		if len(digits) != 12 || !strings.HasPrefix(digits, "98") {
			return fmt.Errorf("invalid mobile format")
		}
		local := digits[2:]
		if len(local) != 10 || local[0] != '9' {
			return fmt.Errorf("invalid mobile format")
		}
		return nil
	}

	if !strings.HasPrefix(digits, ccDigits) {
		return fmt.Errorf("invalid mobile format")
	}
	if len(digits) <= len(ccDigits) {
		return fmt.Errorf("invalid mobile format")
	}
	return nil
}
