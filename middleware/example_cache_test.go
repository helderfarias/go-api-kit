package middleware

import (
	"context"
	"fmt"

	"github.com/helderfarias/go-api-kit/cache"
	"github.com/helderfarias/go-api-kit/endpoint"
	"github.com/spf13/viper"
)

func ExampleCacheServer() {
	viper.Set("cache_redis_servers", "localhost:6379,0,senha,false,*.localhost")
	cacheServer := cache.NewCacheServer()

	service := func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "results"), nil
	}
	service = CachePut(cacheServer, "addresses")(service)
	service = Cacheable(cacheServer, "addresses")(service)
	service = CacheEvict(cacheServer, "addresses")(service)

	resp, err := service(nil, "request")
	fmt.Println(resp.Data())
	fmt.Println(err)

	// Output
	// results
	// nil
}
