package impl

import (
	"cache-service/internal/cachestore"
	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/otp/service"
	"context"
	"errors"
)

// cacheService implements CacheService interface
type cacheService struct {
	store *cachestore.Cache[string, string]
}

// NewCacheService creates a new cache service
func NewOtpCacheService(store *cachestore.Cache[string, string]) service.OtpCacheService {
	return &cacheService{store: store}
}

// CreateCache creates a new cache with validation
func (s *cacheService) AddItem(ctx context.Context, key string, otp string) error {
	err := s.store.Set(ctx, key, otp)

	if err != nil {
		var ec *cachestore.ErrorCode
		if errors.As(err, &ec) {
			return exception.NewAppError(
				ec.Code, ec.HttpStatus, key, otp).WithCause(err)
		}
	}

	return nil
}

func (s *cacheService) GetCacheByKey(ctx context.Context, key string) (string, error) {
	item, err := s.store.GetItem(ctx, key)

	if err != nil {
		var ec *cachestore.ErrorCode
		if errors.As(err, &ec) {
			return "", exception.NewAppError(
				ec.Code, ec.HttpStatus, key).WithCause(err)
		}

		return "", err
	}
	return item.Value, nil
}
