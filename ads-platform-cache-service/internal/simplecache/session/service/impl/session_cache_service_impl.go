package impl

import (
	"context"
	"errors"

	"cache-service/internal/cachestore"
	"cache-service/internal/core/exception"
	"cache-service/internal/simplecache/session/service"
)

type sessionCacheService struct {
	store *cachestore.Cache[string, string]
}

func NewSessionCacheService(store *cachestore.Cache[string, string]) service.SessionCacheService {
	return &sessionCacheService{store: store}
}

func (s *sessionCacheService) AddItem(ctx context.Context, key string, data string) error {
	if err := s.store.Set(ctx, key, data); err != nil {
		var ec *cachestore.ErrorCode
		if errors.As(err, &ec) {
			return exception.NewAppError(ec.Code, ec.HttpStatus, key, data).WithCause(err)
		}
		return err
	}
	return nil
}

func (s *sessionCacheService) GetCacheByKey(ctx context.Context, key string) (string, error) {
	item, err := s.store.GetItem(ctx, key)
	if err != nil {
		var ec *cachestore.ErrorCode
		if errors.As(err, &ec) {
			return "", exception.NewAppError(ec.Code, ec.HttpStatus, key).WithCause(err)
		}
		return "", err
	}
	return item.Value, nil
}

func (s *sessionCacheService) DeleteItem(ctx context.Context, key string) error {
	return s.store.Delete(ctx, key)
}
