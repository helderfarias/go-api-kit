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

type CacheOptions struct {
	TTL time.Duration
}

// Cacheable The simplest way to enable caching behavior for a method is to demarcate it
// with Cacheable and parameterize it with the name of the cache where the results would be stored
func Cacheable(cache cache.CacheServer, name string, options ...CacheOptions) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(parent context.Context, request interface{}) (endpoint.EndpointResponse, error) {
			var target endpoint.EndpointResponse

			key := keyGenerador(name, request)

			cached, err := cache.Get(key, &target)
			if err != nil {
				logrus.Error(err)
			} else if cached != nil {
				return cached.(endpoint.EndpointResponse), nil
			}

			resp, err := next(parent, request)

			if err == nil && resp != nil {
				defaultOptions := CacheOptions{TTL: time.Duration(0)}
				if len(options) >= 1 {
					defaultOptions = options[0]
				}

				if err := cache.Set(key, resp, defaultOptions.TTL); err != nil {
					logrus.Error(err)
				}
			}

			return resp, err
		}
	}
}

func keyGenerador(name, args interface{}) string {
	algorithm := md5.New()
	algorithm.Write([]byte(fmt.Sprintf("%v.%v", name, args)))
	return hex.EncodeToString(algorithm.Sum(nil))
}
