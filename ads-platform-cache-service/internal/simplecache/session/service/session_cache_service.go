package service

import "context"

type SessionCacheService interface {
	AddItem(ctx context.Context, key string, data string) error
	GetCacheByKey(ctx context.Context, key string) (string, error)
	DeleteItem(ctx context.Context, key string) error
}
