package cachestore

import "time"

const (
	DefaultTTL          = 5 * time.Minute
	DefaultCleanupEvery = 10 * time.Minute
	NoExpiration        = time.Duration(-1)
)

type Options struct {
	DefaultTTL      time.Duration
	CleanupInterval time.Duration
	MaxItems        int
	OnEvict         func(key string)
}

func defaultOptions() Options {
	return Options{
		DefaultTTL:      DefaultTTL,
		CleanupInterval: DefaultCleanupEvery,
		MaxItems:        0,
		OnEvict:         nil,
	}
}

type OptionFunc func(*Options)

func WithTTL(ttl time.Duration) OptionFunc {
	return func(o *Options) {
		o.DefaultTTL = ttl
	}
}

func WithCleanupInterval(d time.Duration) OptionFunc {
	return func(o *Options) {
		o.CleanupInterval = d
	}
}

func WithMaxItems(max int) OptionFunc {
	return func(o *Options) {
		o.MaxItems = max
	}
}

func WithOnEvict(fn func(key string)) OptionFunc {
	return func(o *Options) {
		o.OnEvict = fn
	}
}
