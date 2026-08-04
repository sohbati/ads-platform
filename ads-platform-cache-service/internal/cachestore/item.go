package cachestore

import "time"

type Item[T any] struct {
	Value     T
	CreatedAt time.Time
	ExpiresAt time.Time
	HitCount  int
}

func NewItem[T any](value T, ttl time.Duration) Item[T] {
	now := time.Now()
	return Item[T]{
		Value:     value,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
		HitCount:  0,
	}
}

func (i *Item[T]) IsExpired() bool {
	if i.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(i.ExpiresAt)
}

func (i *Item[T]) TTLRemaining() time.Duration {
	if i.ExpiresAt.IsZero() {
		return -1
	}
	remaining := time.Until(i.ExpiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}
