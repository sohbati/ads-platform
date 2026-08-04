package cachestore

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Stats holds cache metrics
type Stats struct {
	Hits      int64
	Misses    int64
	Evictions int64
	ItemCount int64
}

// Cache is a generic thread-safe in-memory cache
type Cache[K comparable, V any] struct {
	mu     sync.RWMutex
	items  map[K]Item[V]
	opts   Options
	stopCh chan struct{}
	stats  Stats
}

// New creates a new generic cache
func New[K comparable, V any](optFuncs ...OptionFunc) *Cache[K, V] {
	opts := defaultOptions()
	for _, fn := range optFuncs {
		fn(&opts)
	}

	c := &Cache[K, V]{
		items:  make(map[K]Item[V]),
		opts:   opts,
		stopCh: make(chan struct{}),
	}

	if opts.CleanupInterval > 0 {
		go c.cleanupLoop()
	}

	return c
}

// ─────────────────────────────────────────────
//  WRITE OPERATIONS
// ─────────────────────────────────────────────

// Set stores a value with the default TTL
func (c *Cache[K, V]) Set(ctx context.Context, key K, value V) error {
	return c.SetWithTTL(ctx, key, value, c.opts.DefaultTTL)
}

// SetWithTTL stores a value with a custom TTL
func (c *Cache[K, V]) SetWithTTL(ctx context.Context, key K, value V, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// enforce max items
	if c.opts.MaxItems > 0 && len(c.items) >= c.opts.MaxItems {
		c.evictOldest()
	}

	var item Item[V]
	if ttl == NoExpiration {
		item = Item[V]{
			Value:     value,
			CreatedAt: time.Now(),
			ExpiresAt: time.Time{}, // zero means never expires
		}
	} else {
		item = NewItem(value, ttl)
	}

	c.items[key] = item
	c.stats.ItemCount = int64(len(c.items))
	return nil
}

// ─────────────────────────────────────────────
//  READ OPERATIONS
// ─────────────────────────────────────────────

// Get retrieves a value by key
func (c *Cache[K, V]) Get(ctx context.Context, key K) (V, error) {
	select {
	case <-ctx.Done():
		var zero V
		return zero, ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		c.stats.Misses++
		var zero V
		return zero, ErrCacheMiss
	}

	if item.IsExpired() {
		delete(c.items, key)
		c.stats.Misses++
		c.stats.Evictions++
		c.stats.ItemCount = int64(len(c.items))
		var zero V
		return zero, ErrCacheMiss
	}

	item.HitCount++
	c.items[key] = item
	c.stats.Hits++
	return item.Value, nil
}

// GetItem returns full Item with metadata
func (c *Cache[K, V]) GetItem(ctx context.Context, key K) (Item[V], error) {
	select {
	case <-ctx.Done():
		return Item[V]{}, ctx.Err()
	default:
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return Item[V]{}, ErrCacheMiss
	}
	if item.IsExpired() {
		return Item[V]{}, ErrKeyExpired
	}
	return item, nil
}

// GetOrSet returns existing value or calls loader to set a new one
func (c *Cache[K, V]) GetOrSet(ctx context.Context, key K, loader func() (V, error)) (V, error) {
	val, err := c.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	newVal, err := loader()
	if err != nil {
		var zero V
		return zero, err
	}

	if setErr := c.Set(ctx, key, newVal); setErr != nil {
		var zero V
		return zero, setErr
	}

	return newVal, nil
}

// GetOrSetWithTTL like GetOrSet but with a custom TTL
func (c *Cache[K, V]) GetOrSetWithTTL(
	ctx context.Context,
	key K,
	ttl time.Duration,
	loader func() (V, error),
) (V, error) {
	val, err := c.Get(ctx, key)
	if err == nil {
		return val, nil
	}

	newVal, err := loader()
	if err != nil {
		var zero V
		return zero, err
	}

	if setErr := c.SetWithTTL(ctx, key, newVal, ttl); setErr != nil {
		var zero V
		return zero, setErr
	}

	return newVal, nil
}

// ─────────────────────────────────────────────
//  DELETE OPERATIONS
// ─────────────────────────────────────────────

// Delete removes a single key
func (c *Cache[K, V]) Delete(ctx context.Context, key K) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	c.stats.ItemCount = int64(len(c.items))
	return nil
}

// DeleteMany removes multiple keys at once
func (c *Cache[K, V]) DeleteMany(ctx context.Context, keys ...K) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		delete(c.items, key)
	}
	c.stats.ItemCount = int64(len(c.items))
	return nil
}

// Flush removes all items
func (c *Cache[K, V]) Flush(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]Item[V])
	c.stats.ItemCount = 0
	return nil
}

// ─────────────────────────────────────────────
//  UTILITY
// ─────────────────────────────────────────────

// Exists returns true if key exists and is not expired
func (c *Cache[K, V]) Exists(ctx context.Context, key K) bool {
	_, err := c.Get(ctx, key)
	return err == nil
}

// Count returns the number of items in cache (including not-yet-cleaned expired ones)
func (c *Cache[K, V]) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Keys returns all non-expired keys
func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.items))
	for k, v := range c.items {
		if !v.IsExpired() {
			keys = append(keys, k)
		}
	}
	return keys
}

// Stats returns a snapshot of cache statistics
func (c *Cache[K, V]) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats
}

// Stop shuts down the background cleanup goroutine
func (c *Cache[K, V]) Stop() {
	close(c.stopCh)
}

// ─────────────────────────────────────────────
//  INTERNAL
// ─────────────────────────────────────────────

func (c *Cache[K, V]) cleanupLoop() {
	ticker := time.NewTicker(c.opts.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopCh:
			return
		}
	}
}

func (c *Cache[K, V]) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.IsExpired() {
			delete(c.items, k)
			c.stats.Evictions++
			if c.opts.OnEvict != nil {
				key := fmt.Sprintf("%v", k)
				go c.opts.OnEvict(key)
			}
		}
	}
	c.stats.ItemCount = int64(len(c.items))
}

func (c *Cache[K, V]) evictOldest() {
	var (
		oldestKey K
		oldest    time.Time
		found     bool
	)

	for k, v := range c.items {
		if !found || v.CreatedAt.Before(oldest) {
			oldest = v.CreatedAt
			oldestKey = k
			found = true
		}
	}

	if found {
		delete(c.items, oldestKey)
		c.stats.Evictions++
		if c.opts.OnEvict != nil {
			key := fmt.Sprintf("%v", oldestKey)
			go c.opts.OnEvict(key)
		}
	}
}
