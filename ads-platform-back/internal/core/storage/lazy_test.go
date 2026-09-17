package storage

import (
	"context"
	"testing"
)

func TestLazyMinioUnavailableWithoutConfig(t *testing.T) {
	s := NewLazyMinio("", "", "", "", "", false)
	if s.Available(context.Background()) {
		t.Fatal("expected unavailable when minio is not configured")
	}
}
