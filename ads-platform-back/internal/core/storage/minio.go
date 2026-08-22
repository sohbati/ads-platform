package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioStorage struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

func NewMinio(endpoint, accessKey, secretKey, bucket, publicURL string, useSSL bool) (ObjectStorage, error) {
	endpoint = strings.TrimSpace(endpoint)
	bucket = strings.TrimSpace(bucket)
	if endpoint == "" || bucket == "" {
		return nil, fmt.Errorf("minio: endpoint and bucket are required")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: connect %s: %w", endpoint, err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("minio: bucket %s: %w", bucket, err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio: create bucket %s: %w", bucket, err)
		}
	}

	publicURL = strings.TrimRight(strings.TrimSpace(publicURL), "/")
	if publicURL == "" {
		scheme := "http"
		if useSSL {
			scheme = "https"
		}
		publicURL = scheme + "://" + endpoint
	}

	return &minioStorage{client: client, bucket: bucket, publicURL: publicURL}, nil
}

func (s *minioStorage) Put(ctx context.Context, key, contentType string, body io.Reader, size int64) (string, error) {
	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}
	if _, err := s.client.PutObject(ctx, s.bucket, key, body, size, opts); err != nil {
		return "", fmt.Errorf("minio: put %s: %w", key, err)
	}
	return s.publicURL + "/" + s.bucket + "/" + key, nil
}
