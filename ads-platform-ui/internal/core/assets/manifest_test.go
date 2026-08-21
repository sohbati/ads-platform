package assets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManifestVersionsTrackContent(t *testing.T) {
	root := t.TempDir()
	css := filepath.Join(root, "css")
	if err := os.Mkdir(css, 0o755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(css, "main.css")
	if err := os.WriteFile(file, []byte("body{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	m, err := Load(root, "/static")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	u := m.URL("css/main.css")
	if !strings.HasPrefix(u, "/static/css/main.css?v=") {
		t.Fatalf("URL = %q, want versioned /static/css/main.css", u)
	}

	// Same content must give a stable version; changed content a new one.
	m2, _ := Load(root, "/static")
	if m2.URL("css/main.css") != u {
		t.Errorf("version changed for identical content")
	}
	if err := os.WriteFile(file, []byte("body{margin:0}"), 0o644); err != nil {
		t.Fatal(err)
	}
	m3, _ := Load(root, "/static")
	if m3.URL("css/main.css") == u {
		t.Errorf("version did not change after content change")
	}

	// Unknown files fall back to an unversioned URL.
	if got := m.URL("js/missing.js"); got != "/static/js/missing.js" {
		t.Errorf("URL(missing) = %q, want unversioned fallback", got)
	}
}
