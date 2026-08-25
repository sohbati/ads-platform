package jsonstore

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// MapStore loads and caches a JSON object file (string keys) with a short TTL.
type MapStore[T any] struct {
	mu       sync.RWMutex
	data     map[string]T
	loadedAt time.Time
	ttl      time.Duration
	path     string
}

func NewMap[T any](path string, ttl time.Duration) *MapStore[T] {
	if ttl <= 0 {
		ttl = 30 * time.Second
	}
	return &MapStore[T]{path: path, ttl: ttl}
}

func (s *MapStore[T]) Get() (map[string]T, error) {
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
	var items map[string]T
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("jsonstore: parse %s: %w", s.path, err)
	}
	if items == nil {
		items = map[string]T{}
	}
	s.data = items
	s.loadedAt = time.Now()
	return s.data, nil
}
