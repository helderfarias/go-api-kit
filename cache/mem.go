package cache

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
)

type memoryCache struct {
	delegate *cache.Cache
}

type memoryCacheValue struct {
	body interface{}
}

func newMemoryCache() *memoryCache {
	delegate := cache.New(5*time.Minute, 10*time.Minute)
	return &memoryCache{delegate: delegate}
}

func (r *memoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	enc, err := json.Marshal(value)
	if err != nil {
		return err
	}

	r.delegate.Set(key, string(enc), ttl)
	return nil
}

func (r *memoryCache) Get(key string, target interface{}) (interface{}, error) {
	dec, ok := r.delegate.Get(key)
	if !ok {
		return "", nil
	}

	return target, json.Unmarshal([]byte(dec.(string)), &target)
}

func (r *memoryCache) Close() error {
	return nil
}
