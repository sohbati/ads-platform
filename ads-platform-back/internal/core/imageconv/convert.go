package imageconv

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
	xwebp "golang.org/x/image/webp"
)

const (
	FullMaxEdge   = 1600
	ThumbMaxEdge  = 480
	FullQuality   = 80
	ThumbQuality  = 75
)

// Variants is the WebP pair stored in object storage.
type Variants struct {
	Full  []byte
	Thumb []byte
}

// FromUpload decodes a JPEG, PNG, or WebP upload and returns a resized
// full-size WebP plus a card thumbnail. The original file is not kept.
func FromUpload(r io.Reader, contentType string) (Variants, error) {
	img, err := decode(r, contentType)
	if err != nil {
		return Variants{}, fmt.Errorf("decode: %w", err)
	}
	if img.Bounds().Empty() {
		return Variants{}, fmt.Errorf("decode: empty image")
	}

	fullImg := downscale(img, FullMaxEdge)
	thumbImg := downscale(img, ThumbMaxEdge)

	full, err := encodeWebP(imaging.Clone(fullImg), FullQuality)
	if err != nil {
		return Variants{}, fmt.Errorf("encode full: %w", err)
	}
	thumb, err := encodeWebP(imaging.Clone(thumbImg), ThumbQuality)
	if err != nil {
		return Variants{}, fmt.Errorf("encode thumb: %w", err)
	}
	return Variants{Full: full, Thumb: thumb}, nil
}

func decode(r io.Reader, contentType string) (image.Image, error) {
	switch contentType {
	case "image/webp":
		return xwebp.Decode(r)
	case "image/png":
		return png.Decode(r)
	case "image/jpeg":
		return jpeg.Decode(r)
	default:
		img, _, err := image.Decode(r)
		return img, err
	}
}

func downscale(img image.Image, maxEdge int) image.Image {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= maxEdge && h <= maxEdge {
		return img
	}
	return imaging.Fit(img, maxEdge, maxEdge, imaging.Lanczos)
}

// PeekWebP reports whether buf looks like a RIFF/WEBP file.
func PeekWebP(buf []byte) bool {
	return len(buf) >= 12 && bytes.Equal(buf[:4], []byte("RIFF")) && bytes.Equal(buf[8:12], []byte("WEBP"))
}
