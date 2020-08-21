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
	s1 := newCacheServer("localhost:8000")
	s2 := newCacheServer(":8000,localhost:8001")

	assert.IsType(t, s1, &redisCache{})
	assert.IsType(t, s2, &redisCache{})
}

func TestCreateServerWithRedisCache(t *testing.T) {
	viper.Set("cache_redis_servers", "localhost:6380,localhost:6380")
	viper.Set("cache_redis_password", "password00")
	viper.Set("cache_redis_ssl", "true")

	s := NewCacheServer()

	assert.NotNil(t, s)
}
