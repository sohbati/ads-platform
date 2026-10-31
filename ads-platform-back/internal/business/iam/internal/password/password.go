package password

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinLength = 8
	// bcrypt ignores everything after 72 bytes, so longer input is rejected rather than truncated.
	MaxBytes = 72
	cost     = 12
)

// dummyHash lets Verify spend bcrypt time even when there is no stored hash, so callers do not leak
// whether an account exists through response timing.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("timing-equalizer-not-a-real-password"), cost)

func Valid(raw string) bool {
	return utf8.RuneCountInString(raw) >= MinLength && len(raw) <= MaxBytes
}

func Hash(raw string) (string, error) {
	if !Valid(raw) {
		return "", errors.New("password does not meet policy")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Verify reports whether raw matches hash. An empty hash never matches.
func Verify(hash, raw string) bool {
	if hash == "" {
		_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(raw))
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw)) == nil
}
