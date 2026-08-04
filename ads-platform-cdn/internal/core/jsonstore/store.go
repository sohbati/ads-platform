package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Store loads and caches a JSON array file with a short TTL.
type Store[T any] struct {
	mu       sync.RWMutex
	data     []T
	loadedAt time.Time
	ttl      time.Duration
	path     string
}

func New[T any](path string, ttl time.Duration) *Store[T] {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &Store[T]{path: path, ttl: ttl}
}

func (s *Store[T]) Get() ([]T, error) {
	s.mu.RLock()
	if s.data != nil && time.Since(s.loadedAt) < s.ttl {
		defer s.mu.RUnlock()
		return s.data, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data != nil && time.Since(s.loadedAt) < s.ttl {
		return s.data, nil
	}

	raw, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("jsonstore: read %s: %w", s.path, err)
	}
	var items []T
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("jsonstore: parse %s: %w", s.path, err)
	}
	s.data = items
	s.loadedAt = time.Now()
	return s.data, nil
}
