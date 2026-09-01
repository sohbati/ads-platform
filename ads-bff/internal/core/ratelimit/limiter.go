package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	window time.Duration
	max    int
}

func New(window time.Duration, max int) *Limiter {
	if max < 1 {
		max = 60
	}
	if window <= 0 {
		window = time.Minute
	}
	return &Limiter{
		hits:   map[string][]time.Time{},
		window: window,
		max:    max,
	}
}

func (l *Limiter) Allow(key string) bool {
	now := time.Now()
	cutoff := now.Add(-l.window)

	l.mu.Lock()
	defer l.mu.Unlock()

	q := l.hits[key]
	i := 0
	for i < len(q) && q[i].Before(cutoff) {
		i++
	}
	q = q[i:]
	if len(q) >= l.max {
		l.hits[key] = q
		return false
	}
	l.hits[key] = append(q, now)
	return true
}
