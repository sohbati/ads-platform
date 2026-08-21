package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Manifest maps static files (relative to the static root) to a short content
// hash used as a cache-busting version. Assets are served with an immutable
// Cache-Control header, so any content change must produce a new URL; hashing
// at startup guarantees that without manual version bumps.
type Manifest struct {
	urlPrefix string
	versions  map[string]string
}

// Load walks root and hashes every file. urlPrefix is the public mount point,
// e.g. "/static".
func Load(root, urlPrefix string) (*Manifest, error) {
	m := &Manifest{urlPrefix: urlPrefix, versions: make(map[string]string)}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || d.Name() == ".DS_Store" {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		version, err := hashFile(path)
		if err != nil {
			return err
		}
		m.versions[filepath.ToSlash(rel)] = version
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("assets: load %s: %w", root, err)
	}
	return m, nil
}

// URL returns the versioned public URL for an asset:
// "css/layout.css" -> "/static/css/layout.css?v=1a2b3c4d".
// Unknown files are returned unversioned so a missing hash never breaks a page.
func (m *Manifest) URL(rel string) string {
	u := m.urlPrefix + "/" + rel
	if v, ok := m.versions[rel]; ok {
		return u + "?v=" + v
	}
	return u
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil))[:8], nil
}
