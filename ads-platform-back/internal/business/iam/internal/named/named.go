package named

import (
	"regexp"
	"strings"
)

var CodePattern = regexp.MustCompile(`^[a-z][a-z0-9._-]{0,79}$`)

func Code(raw string) (string, bool) {
	code := strings.ToLower(strings.TrimSpace(raw))
	if !CodePattern.MatchString(code) {
		return "", false
	}
	return code, true
}

func Name(raw string) (string, bool) {
	name := strings.TrimSpace(raw)
	if name == "" || len(name) > 120 {
		return "", false
	}
	return name, true
}
