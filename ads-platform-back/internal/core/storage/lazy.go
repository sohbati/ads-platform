package storage

import (
	"context"
	"io"
	"sync"
)

// lazyMinio connects on first use and reconnects after MinIO comes back.
type lazyMinio struct {
	mu        sync.Mutex
	inner     ObjectStorage
	endpoint  string
	accessKey string
	secretKey string
	bucket    string
	publicURL string
	useSSL    bool
}

func NewLazyMinio(endpoint, accessKey, secretKey, bucket, publicURL string, useSSL bool) ObjectStorage {
	return &lazyMinio{
		endpoint:  endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		publicURL: publicURL,
		useSSL:    useSSL,
	}
}

func (s *lazyMinio) Put(ctx context.Context, key, contentType string, body io.Reader, size int64) (string, error) {
	inner, err := s.ensure(ctx)
	if err != nil {
		return "", err
	}
	return inner.Put(ctx, key, contentType, body, size)
}

func (s *lazyMinio) Available(ctx context.Context) bool {
	inner, err := s.ensure(ctx)
	return err == nil && inner != nil
}

func (s *lazyMinio) ensure(ctx context.Context) (ObjectStorage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.inner != nil {
		if s.inner.Available(ctx) {
			return s.inner, nil
		}
		s.inner = nil
	}

	inner, err := newMinio(ctx, s.endpoint, s.accessKey, s.secretKey, s.bucket, s.publicURL, s.useSSL)
	if err != nil {
		return nil, err
	}
	s.inner = inner
	return inner, nil
}
