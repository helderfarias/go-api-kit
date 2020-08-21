package cache

import (
	"crypto/tls"
	"encoding/json"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/redis.v5"
)

type redisCache struct {
	master *redis.Client
	slave  *redis.Client
}

func newRedisCache(masterAddr string, slaveAddr string, password string, ssl bool) CacheServer {
	instance := &redisCache{}

	tlsMaster := instance.toTLS(masterAddr, ssl)
	instance.master = redis.NewClient(&redis.Options{Addr: masterAddr, Password: password, TLSConfig: tlsMaster})
	if ok := instance.master.Ping(); ok.Err() != nil {
		logrus.Errorf("Could not connect Redis Master, %v", ok.Err())
		return instance
	}

	if slaveAddr != "" {
		tlsSlave := instance.toTLS(slaveAddr, ssl)
		instance.slave = redis.NewClient(&redis.Options{Addr: slaveAddr, Password: password, TLSConfig: tlsSlave})
		if ok := instance.slave.Ping(); ok.Err() != nil {
			logrus.Errorf("Could not connect Redis Slave, %v", ok.Err())
			return instance
		}
	}

	return instance
}

func (r *redisCache) toTLS(url string, ssl bool) *tls.Config {
	var tlsName *tls.Config

	if ssl {
		serverName := ""
		if len(strings.Split(url, ":")) >= 0 {
			serverName = strings.Split(url, ":")[0]
		}

		tlsName = &tls.Config{ServerName: serverName}
	}

	return tlsName
}

func (r *redisCache) Expire(key string, ttl time.Duration) error {
	if err := r.master.Expire(key, ttl); err != nil {
		logrus.Error(err)
	}

	if err := r.slave.Expire(key, ttl); err != nil {
		logrus.Error(err)
	}

	return nil
}

func (r *redisCache) Close() error {
	if err := r.master.Close(); err != nil {
		logrus.Error(err)
	}

	if err := r.slave.Close(); err != nil {
		logrus.Error(err)
	}

	return nil
}

func (r *redisCache) Set(key string, value interface{}, ttl time.Duration) error {
	enc, err := json.Marshal(value)
	if err != nil {
		return err
	}

	status := r.master.Set(key, enc, ttl)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *redisCache) Get(key string, target interface{}) (interface{}, error) {
	status := r.slave.Get(key)
	if status.Err() != nil {
		return "", status.Err()
	}

	return target, json.Unmarshal([]byte(status.Val()), &target)
}
