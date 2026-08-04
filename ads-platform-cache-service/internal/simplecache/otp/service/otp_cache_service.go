package service

import (
	"context"
)

// CacheService defines the interface for cache business logic operations
type OtpCacheService interface {
	AddItem(ctx context.Context, key string, otp string) error
	GetCacheByKey(ctx context.Context, key string) (string, error)
}
