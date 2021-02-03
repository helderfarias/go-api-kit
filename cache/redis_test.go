package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisCacheWithConfig(t *testing.T) {
	r := newRedisCache("localhost:6380,0,pass123,true,*.redis.localhost", false)

	assert.NotNil(t, r)
	assert.Equal(t, "localhost:6380", r.options.Addr)
	assert.Equal(t, 0, r.options.DB)
	assert.Equal(t, "pass123", r.options.Password)
	assert.Equal(t, "*.redis.localhost", r.options.TLSConfig.ServerName)
}
