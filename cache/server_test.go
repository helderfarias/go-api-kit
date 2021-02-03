package cache

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateMemoryCache(t *testing.T) {
	s := NewCacheServer()

	assert.IsType(t, s, &memoryCache{})
}

func TestCreateRedisCache(t *testing.T) {
	viper.Set("cache_redis_ping", "false")

	s1 := newCacheServer("localhost:8000")

	assert.IsType(t, s1, &redisCache{})
}

func TestCreateServerRedisWithFallbackToMemory(t *testing.T) {
	viper.Set("cache_redis_ping", "true")
	viper.Set("cache_redis_servers", "fallback:4444")

	s := NewCacheServer()

	assert.IsType(t, &memoryCache{}, s)
}
