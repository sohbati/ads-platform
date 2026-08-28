package storage

import (
	"context"
	"io"
)

// ObjectStorage stores ad pictures in an S3-compatible backend (MinIO).
// Put returns a host-free path: /{bucket}/{key}.
type ObjectStorage interface {
	Put(ctx context.Context, key, contentType string, body io.Reader, size int64) (publicPath string, err error)
}
