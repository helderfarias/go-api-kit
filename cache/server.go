package cache

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CacheServer interface {
	Set(key string, value interface{}, ttl time.Duration) error

	Get(key string, target interface{}) (interface{}, error)

	Expire(key string, ttl time.Duration) error

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

	redisServers := strings.Split(servers, ",")
	if len(redisServers) == 0 {
		logrus.Info("Working with Memory Cache")
		return newMemoryCache()
	}

	logrus.Infof("Working with Redis Cache %d (s)", len(redisServers))

	password := viper.GetString("cache_redis_password")
	ssl := viper.GetBool("cache_redis_ssl")

	if len(redisServers) == 1 {
		return newRedisCache(redisServers[0], "", password, ssl)
	}

	return newRedisCache(redisServers[0], redisServers[1], password, ssl)
}
