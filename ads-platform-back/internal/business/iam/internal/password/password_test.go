package password

import (
	"strings"
	"testing"
)

func TestHashAndVerify(t *testing.T) {
	hash, err := Hash("correct horse")
	if err != nil {
		t.Fatal(err)
	}
	if !Verify(hash, "correct horse") {
		t.Fatal("expected match")
	}
	if Verify(hash, "wrong horse!") {
		t.Fatal("expected mismatch")
	}
}

func TestEmptyHashNeverMatches(t *testing.T) {
	if Verify("", "") || Verify("", "anything") {
		t.Fatal("empty hash matched")
	}
}

func TestPolicy(t *testing.T) {
	if Valid("short") {
		t.Fatal("short password accepted")
	}
	if Valid(strings.Repeat("a", MaxBytes+1)) {
		t.Fatal("over-long password accepted")
	}
	if !Valid("رمز‌عبور‌امن") {
		t.Fatal("persian password rejected")
	}
}
