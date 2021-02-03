package middleware

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/helderfarias/go-api-kit/cache"
	"github.com/helderfarias/go-api-kit/endpoint"
	"github.com/sirupsen/logrus"
)

// CacheEvictOptions cache configurations
type CacheEvictOptions struct {
	TTL          time.Duration
	AllEntries   bool
	OnListener   func(event string, key string)
	KeyGenerator func(name string, args interface{}) string
}

// CacheableOptions cache configurations
type CacheableOptions struct {
	TTL          time.Duration
	OnListener   func(event string, key string)
	KeyGenerator func(name string, args interface{}) string
}

// CachePutOptions cache configurations
type CachePutOptions struct {
	TTL          time.Duration
	OnListener   func(event string, key string)
	KeyGenerator func(name string, args interface{}) string
}

// EntryCache cache container
type entryCache struct {
	Status int         `json:"status"`
	Value  interface{} `json:"value"`
}

// DefaultListener listener
var DefaultListener = func(event string, nameOrKey string) {}

// DefaultKeyGenerator generator
var DefaultKeyGenerator = func(name string, args interface{}) string {
	algorithm := md5.New()
	algorithm.Write([]byte(fmt.Sprintf("%v%v", name, args)))
	return fmt.Sprintf("%v:%v", name, hex.EncodeToString(algorithm.Sum(nil)))
}

// CacheEvict Now, what would be the problem with making all methods Cacheable?
// The problem is size – we don't want to populate the cache with values that we don't need often.
// Caches can grow quite large, quite fast, and we could be holding on to a lot of stale or unused data.
// The CacheEvict is used to indicate the removal of one or more/all values – so that fresh values can be loaded into the cache again.
func CacheEvict(cache cache.CacheServer, name string, options ...CacheEvictOptions) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(parent context.Context, request interface{}) (endpoint.EndpointResponse, error) {
			opt := CacheEvictOptions{TTL: time.Duration(0), OnListener: DefaultListener, KeyGenerator: DefaultKeyGenerator}
			if len(options) >= 1 {
				opt = options[0]
				if opt.OnListener == nil {
					opt.OnListener = DefaultListener
				}
				if opt.KeyGenerator == nil {
					opt.KeyGenerator = DefaultKeyGenerator
				}
			}

			if opt.AllEntries {
				if err := cache.DeleteAll(name); err != nil {
					logrus.Error(err)
				} else {
					logrus.WithField("cacheable.clean.allentries", "true").Debug("CacheEvict")
					opt.OnListener("evict", name)
				}
			} else {
				if err := cache.Delete(name); err != nil {
					logrus.Error(err)
				} else {
					logrus.WithField("cacheable.clean.onlyentry", "true").Debug("CacheEvict")
					opt.OnListener("evict", name)
				}
			}

			return next(parent, request)
		}
	}
}

// CachePut While CacheEvict reduces the overhead of looking up entries in a large cache by removing stale and unused entries,
// ideally, you want to avoid evicting too much data out of the cache.
// Instead, you'd want to selectively and intelligently update the entries whenever they're altered.
// With the CachePut annotation, you can update the content of the cache without interfering the method execution.
// That is, the method would always be executed and the result cached.
func CachePut(cache cache.CacheServer, name string, options ...CachePutOptions) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(parent context.Context, request interface{}) (endpoint.EndpointResponse, error) {
			opt := CachePutOptions{TTL: time.Duration(0), OnListener: DefaultListener, KeyGenerator: DefaultKeyGenerator}
			if len(options) >= 1 {
				opt = options[0]
				if opt.OnListener == nil {
					opt.OnListener = DefaultListener
				}
				if opt.KeyGenerator == nil {
					opt.KeyGenerator = DefaultKeyGenerator
				}
			}

			if err := cache.Delete(name); err != nil {
				logrus.Error(err)
			} else {
				logrus.WithField("cacheable.clean.onlyentry", "true").Debug("CachePut")
				opt.OnListener("evict", name)
			}

			resp, err := next(parent, request)

			if err == nil && resp.Data() != nil {
				key := opt.KeyGenerator(name, request)

				newEntry := entryCache{Status: resp.Code(), Value: resp.Data()}

				if err := cache.Set(key, newEntry, opt.TTL); err != nil {
					logrus.Error(err)
				} else {
					logrus.WithField("cacheable.put", key).Debug("CachePut")
					opt.OnListener("put", key)
				}
			}

			return resp, err
		}
	}
}

// Cacheable The simplest way to enable caching behavior for a method is to demarcate it
// with Cacheable and parameterize it with the name of the cache where the results would be stored
func Cacheable(cache cache.CacheServer, name string, options ...CacheableOptions) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(parent context.Context, request interface{}) (endpoint.EndpointResponse, error) {
			opt := CacheableOptions{TTL: time.Duration(0), OnListener: DefaultListener, KeyGenerator: DefaultKeyGenerator}
			if len(options) >= 1 {
				opt = options[0]
				if opt.OnListener == nil {
					opt.OnListener = DefaultListener
				}
				if opt.KeyGenerator == nil {
					opt.KeyGenerator = DefaultKeyGenerator
				}
			}

			key := opt.KeyGenerator(name, request)

			var entry entryCache
			cached, err := cache.Get(key, &entry)
			if err != nil {
				logrus.Error(err)

				if err := cache.Delete(key); err != nil {
					logrus.Error(err)
				}
			} else if cache, ok := cached.(*entryCache); ok && cache.Value != nil {
				logrus.WithField("cacheable.get", key).Debug("Cacheable")
				opt.OnListener("get", key)
				return endpoint.Response(cache.Status, cache.Value), nil
			}

			resp, err := next(parent, request)

			if err == nil && resp.Data() != nil {
				newEntry := entryCache{Status: resp.Code(), Value: resp.Data()}

				if err := cache.Set(key, newEntry, opt.TTL); err != nil {
					logrus.Error(err)
				} else {
					logrus.WithField("cacheable.put", key).Debug("Cacheable")
					opt.OnListener("put", key)
				}
			}

			return resp, err
		}
	}
}
