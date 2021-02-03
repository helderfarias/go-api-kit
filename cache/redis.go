package cache

import (
	"crypto/tls"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/redis.v5"
)

type redisCache struct {
	redis   *redis.Client
	options redis.Options
}

func newRedisCache(server string, ping bool) *redisCache {
	options := strings.Split(server, ",")

	redisOpt := redis.Options{
		Addr: options[0],
	}

	if len(options) > 2 {
		if db, err := strconv.Atoi(options[1]); err == nil {
			redisOpt.DB = db
		}
	}

	if len(options) > 3 {
		redisOpt.Password = options[2]
	}

	if len(options) > 4 {
		if ssl, err := strconv.ParseBool(options[3]); err == nil && ssl {
			redisOpt.TLSConfig = buildTLS(options[4])
		}
	}

	redis := redis.NewClient(&redisOpt)

	if ping {
		if ok := redis.Ping(); ok.Err() != nil {
			logrus.Errorf("Could not connect Redis Master, %v", ok.Err())
			return nil
		}
	}

	return &redisCache{redis: redis, options: redisOpt}
}

func (r *redisCache) Expire(key string, ttl time.Duration) error {
	if r.redis == nil {
		logrus.Info("Redis Master is not configured")
		return nil
	}

	if err := r.redis.Expire(key, ttl); err != nil {
		logrus.Error(err)
	}

	return nil
}

func (r *redisCache) Close() error {
	if r.redis == nil {
		logrus.Info("Redis Master is not configured")
		return nil
	}

	if err := r.redis.Close(); err != nil {
		logrus.Error(err)
	}

	return nil
}

func (r *redisCache) Set(key string, value interface{}, ttl time.Duration) error {
	if r.redis == nil {
		logrus.Info("Redis Master is not configured")
		return nil
	}

	enc, err := json.Marshal(value)
	if err != nil {
		return err
	}

	status := r.redis.Set(key, enc, ttl)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *redisCache) Get(key string, target interface{}) (interface{}, error) {
	if r.redis == nil {
		logrus.Info("Redis Master is not configured")
		return nil, nil
	}

	status := r.redis.Get(key)
	if status.Err() != nil {
		return "", status.Err()
	}

	return target, json.Unmarshal([]byte(status.Val()), &target)
}

func buildTLS(serverName string) *tls.Config {
	return &tls.Config{
		ServerName: serverName,
	}
}
