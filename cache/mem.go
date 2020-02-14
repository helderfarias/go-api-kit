package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type memoryCache struct {
	delegate *cache.Cache
}

func newMemoryCache() *memoryCache {
	delegate := cache.New(5*time.Minute, 10*time.Minute)
	return &memoryCache{delegate: delegate}
}

func (r *memoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	r.delegate.Set(key, value, ttl)
	return nil
}

func (r *memoryCache) Get(key string, target interface{}) (interface{}, error) {
	d, ok := r.delegate.Get(key)
	if !ok {
		return "", nil
	}

	return d, nil
}

func (r *memoryCache) Close() error {
	return nil
}
