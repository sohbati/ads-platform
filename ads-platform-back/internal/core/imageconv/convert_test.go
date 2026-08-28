package imageconv

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"testing"
)

func jpeg1x1(t *testing.T) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{R: 200, G: 40, B: 40, A: 255})
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func png1x1(t *testing.T) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.NRGBA{R: 20, G: 180, B: 80, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestFromUploadJPEGProducesWebPPair(t *testing.T) {
	got, err := FromUpload(bytes.NewReader(jpeg1x1(t)), "image/jpeg")
	if err != nil {
		t.Fatal(err)
	}
	if !PeekWebP(got.Full) || !PeekWebP(got.Thumb) {
		t.Fatalf("expected webp riff, full=%d thumb=%d", len(got.Full), len(got.Thumb))
	}
}

func TestFromUploadPNGProducesWebPPair(t *testing.T) {
	got, err := FromUpload(bytes.NewReader(png1x1(t)), "image/png")
	if err != nil {
		t.Fatal(err)
	}
	if !PeekWebP(got.Full) {
		t.Fatal("png upload was not converted to webp")
	}
}

func TestFromUploadWebPRoundTrip(t *testing.T) {
	first, err := FromUpload(bytes.NewReader(jpeg1x1(t)), "image/jpeg")
	if err != nil {
		t.Fatal(err)
	}
	got, err := FromUpload(bytes.NewReader(first.Full), "image/webp")
	if err != nil {
		t.Fatal(err)
	}
	if !PeekWebP(got.Full) || !PeekWebP(got.Thumb) {
		t.Fatal("webp upload was not re-encoded as webp")
	}
}

func TestFromUploadRejectsGarbage(t *testing.T) {
	_, err := FromUpload(bytes.NewReader([]byte("not-an-image")), "image/jpeg")
	if err == nil {
		t.Fatal("expected decode error")
	}
}
