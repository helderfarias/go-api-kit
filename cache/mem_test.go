package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemCreateMemoryCache(t *testing.T) {
	m := newMemoryCache()

	m.Set("key", "value", 10*time.Minute)

	cached := ""
	_, err := m.Get("key", &cached)

	assert.Nil(t, err)
	assert.Equal(t, "value", cached)
}
