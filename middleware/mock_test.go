package middleware

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type cacheServerMock struct {
	mock.Mock
}

func (c *cacheServerMock) Delete(key string) error {
	args := c.Called(key)
	return args.Error(0)
}

func (c *cacheServerMock) DeleteAll(key string) error {
	args := c.Called(key)
	return args.Error(0)
}

func (c *cacheServerMock) Set(key string, value interface{}, ttl time.Duration) error {
	args := c.Called(key, value, ttl)
	return args.Error(0)
}

func (c *cacheServerMock) Get(key string, target interface{}) (interface{}, error) {
	args := c.Called(key, target)
	return args.Get(0), args.Error(1)
}

func (c *cacheServerMock) Expire(key string, ttl time.Duration) error {
	args := c.Called(key, ttl)
	return args.Error(0)
}

func (c *cacheServerMock) Close() error {
	args := c.Called()
	return args.Error(0)
}
