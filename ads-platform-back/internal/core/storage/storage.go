package storage

import (
	"context"
	"io"
)

// ObjectStorage stores ad pictures in an S3-compatible backend (MinIO).
type ObjectStorage interface {
	Put(ctx context.Context, key, contentType string, body io.Reader, size int64) (publicURL string, err error)
}
