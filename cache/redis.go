package cache

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
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

func (r *redisCache) DeleteAll(key string) error {
	if r.redis == nil {
		return errors.New("Redis Master is not configured")
	}

	if strings.TrimSpace(key) == "" {
		return errors.New("Key is empty")
	}

	cmd := r.redis.Keys(fmt.Sprintf("*%v*", key))
	if cmd.Err() != nil {
		return cmd.Err()
	}

	entries, err := cmd.Result()
	if err != nil {
		return err
	}

	var resulErrors error
	for _, key := range entries {
		if cmd := r.redis.Del(key); cmd != nil && cmd.Err() != nil {
			logrus.Error(cmd.Err())
			resulErrors = cmd.Err()
		}
	}

	return resulErrors
}

func (r *redisCache) Delete(key string) error {
	if r.redis == nil {
		return errors.New("Redis Master is not configured")
	}

	if strings.TrimSpace(key) == "" {
		return errors.New("Key is empty")
	}

	cmd := r.redis.Del(key)

	if cmd != nil && cmd.Err() != nil {
		logrus.Error(cmd.Err())
		return cmd.Err()
	}

	return nil
}

func (r *redisCache) Expire(key string, ttl time.Duration) error {
	if r.redis == nil {
		return errors.New("Redis Master is not configured")
	}

	if strings.TrimSpace(key) == "" {
		return errors.New("Key is empty")
	}

	cmd := r.redis.Del(key)
	if cmd != nil && cmd.Err() != nil {
		logrus.Error(cmd.Err())
	}

	return nil
}

func (r *redisCache) Close() error {
	if r.redis == nil {
		return errors.New("Redis Master is not configured")
	}

	if err := r.redis.Close(); err != nil {
		logrus.Error(err)
	}

	return nil
}

func (r *redisCache) Set(key string, value interface{}, ttl time.Duration) error {
	if r.redis == nil {
		return errors.New("Redis Master is not configured")
	}

	if strings.TrimSpace(key) == "" {
		return errors.New("Key is empty")
	}

	if value == nil {
		logrus.Warn("Could not set value is empty for 'Set'")
		return nil
	}

	enc, err := json.Marshal(value)
	if err != nil {
		return err
	}

	cmd := r.redis.Set(key, enc, ttl)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *redisCache) Get(key string, target interface{}) (interface{}, error) {
	if r.redis == nil {
		return nil, errors.New("Redis Master is not configured")
	}

	if strings.TrimSpace(key) == "" {
		return nil, errors.New("Key is empty")
	}

	if target == nil {
		logrus.Warn("Target is nil for 'Get'")
		return target, nil
	}

	cmd := r.redis.Get(key)
	if cmd != nil && cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return target, json.Unmarshal([]byte(cmd.Val()), &target)
}

func buildTLS(serverName string) *tls.Config {
	return &tls.Config{
		ServerName: serverName,
	}
}
