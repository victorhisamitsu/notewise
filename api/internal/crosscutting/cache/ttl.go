// Package cache provides a small in-memory TTL cache with request
// coalescing. On Vercel each function instance keeps its own copy, which
// complements the CDN cache configured through Cache-Control headers.
package cache

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

type entry[V any] struct {
	value   V
	expires time.Time
}

// TTL is a concurrency-safe cache whose entries expire after a fixed time.
// Get-or-load calls for the same key are coalesced with singleflight, so a
// burst of identical requests reaches the upstream API only once.
type TTL[V any] struct {
	mu       sync.RWMutex
	items    map[string]entry[V]
	ttl      time.Duration
	maxItems int
	group    singleflight.Group
	now      func() time.Time
}

// NewTTL creates a cache holding at most maxItems entries for ttl each.
func NewTTL[V any](ttl time.Duration, maxItems int) *TTL[V] {
	return &TTL[V]{items: make(map[string]entry[V]), ttl: ttl, maxItems: maxItems, now: time.Now}
}

// Get returns a fresh cached value.
func (c *TTL[V]) Get(key string) (V, bool) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()
	if !ok || c.now().After(e.expires) {
		var zero V
		return zero, false
	}
	return e.value, true
}

// Set stores a value. When full, expired entries are evicted first and,
// if that is not enough, the cache is reset — simple and bounded.
func (c *TTL[V]) Set(key string, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.items) >= c.maxItems {
		now := c.now()
		for k, e := range c.items {
			if now.After(e.expires) {
				delete(c.items, k)
			}
		}
		if len(c.items) >= c.maxItems {
			c.items = make(map[string]entry[V])
		}
	}
	c.items[key] = entry[V]{value: v, expires: c.now().Add(c.ttl)}
}

// GetOrLoad returns the cached value or calls load once per key.
// Errors are never cached. hit reports whether the cache answered.
func (c *TTL[V]) GetOrLoad(ctx context.Context, key string, load func(context.Context) (V, error)) (v V, hit bool, err error) {
	if v, ok := c.Get(key); ok {
		return v, true, nil
	}
	res, err, _ := c.group.Do(key, func() (any, error) {
		if v, ok := c.Get(key); ok {
			return v, nil
		}
		v, err := load(ctx)
		if err != nil {
			return v, err
		}
		c.Set(key, v)
		return v, nil
	})
	if err != nil {
		var zero V
		return zero, false, err
	}
	return res.(V), false, nil
}
