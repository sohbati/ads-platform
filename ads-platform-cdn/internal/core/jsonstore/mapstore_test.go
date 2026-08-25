package jsonstore

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMapStoreGet(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "attr-schemas.json")
	if err := os.WriteFile(path, []byte(`{"cars":{"title":"خودرو"}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	store := NewMap[struct {
		Title string `json:"title"`
	}](path, time.Minute)

	got, err := store.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got["cars"].Title != "خودرو" {
		t.Fatalf("title: got %q", got["cars"].Title)
	}
}
