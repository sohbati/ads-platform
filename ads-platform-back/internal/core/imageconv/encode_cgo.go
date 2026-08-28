package imageconv

import (
	"fmt"
	"image"
	"unsafe"
)

/*
#cgo linux pkg-config: libwebp
#cgo darwin CFLAGS: -I/opt/homebrew/include -I/usr/local/include
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -lwebp
#include <webp/encode.h>
#include <stdlib.h>
*/
import "C"

func encodeWebP(img *image.NRGBA, quality float32) ([]byte, error) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("empty image")
	}
	if len(img.Pix) == 0 {
		return nil, fmt.Errorf("empty pixels")
	}

	var out *C.uint8_t
	size := C.WebPEncodeRGBA(
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(w),
		C.int(h),
		C.int(img.Stride),
		C.float(quality),
		&out,
	)
	if size == 0 || out == nil {
		return nil, fmt.Errorf("libwebp encode failed")
	}
	defer C.WebPFree(unsafe.Pointer(out))
	return C.GoBytes(unsafe.Pointer(out), C.int(size)), nil
}
