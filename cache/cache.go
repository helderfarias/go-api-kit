package cache

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CacheServer interface {
	Set(key string, value interface{}, ttl time.Duration) error

	Get(key string, target interface{}) (interface{}, error)

	Expire(key string, ttl time.Duration) error

	Delete(key string) error

	DeleteAll(key string) error

	Close() error
}

func NewCacheServer() CacheServer {
	return newCacheServer(viper.GetString("cache_redis_servers"))
}

func newCacheServer(servers string) CacheServer {
	if servers == "" {
		logrus.Info("Working with Memory Cache")
		return newMemoryCache()
	}

	ping := true
	if viper.GetString("cache_redis_ping") == "false" {
		ping = false
	}

	redis := newRedisCache(servers, ping)
	if redis == nil {
		logrus.Infof("Fallback to Memory Cache")
		return newMemoryCache()
	}

	logrus.Infof("Working with Redis Cache")
	return redis
}
