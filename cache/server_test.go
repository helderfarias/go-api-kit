package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMemoryCache(t *testing.T) {
	s := NewCacheServer()

	assert.IsType(t, s, &memoryCache{})
}

func TestCreateRedisCache(t *testing.T) {
	s1 := newCacheServer("localhost:8000")
	s2 := newCacheServer("localhost:8000,localhost:8001")

	assert.IsType(t, s1, &redisCache{})
	assert.IsType(t, s2, &redisCache{})
}
