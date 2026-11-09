package sms

import "strings"

// NormalizeMobile turns an Iranian mobile number in any common format
// (09121234567, 9121234567, 989121234567, +989121234567, 00989121234567,
// with spaces or dashes, Latin or Persian digits) into 09121234567.
// ok is false if it is not a valid mobile number.
func NormalizeMobile(raw string) (normalized string, ok bool) {
	var b strings.Builder
	for i, r := range strings.TrimSpace(raw) {
		switch {
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= '۰' && r <= '۹': // Persian digits
			b.WriteRune('0' + (r - '۰'))
		case r == '+' && i == 0, r == ' ', r == '-':
		default:
			return "", false
		}
	}
	d := b.String()
	switch {
	case strings.HasPrefix(d, "0098"):
		d = d[4:]
	case strings.HasPrefix(d, "98") && len(d) == 12:
		d = d[2:]
	case strings.HasPrefix(d, "0"):
		d = d[1:]
	}
	if len(d) != 10 || d[0] != '9' {
		return "", false
	}
	return "0" + d, true
}
