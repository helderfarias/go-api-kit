package middleware

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/helderfarias/go-api-kit/endpoint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestValueFromCacheWhenNotEmpty(t *testing.T) {
	cacheServerMock := &cacheServerMock{}

	cacheServerMock.On("Get", "7e7440170f67c54ddfdce2035c85d482", mock.Anything).Return(endpoint.Response(200, "cached: address 10"), nil)

	mw := Cacheable(cacheServerMock, "addresses")(func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "not return"), nil
	})

	resp, err := mw(nil, "request")

	assert.Nil(t, err)
	assert.Equal(t, "cached: address 10", resp.Data())
}

func TestPutValueToCacheWhenEmpty(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Get", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything).Return(nil, nil)
	cacheMock.On("Set", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything, 5*time.Minute).Return(nil)

	mw := Cacheable(cacheMock, "addresses", CacheOptions{TTL: 5 * time.Minute})(func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "cached empty"), nil
	})

	resp, err := mw(nil, "params")

	assert.Nil(t, err)
	assert.Equal(t, "cached empty", resp.Data())
	cacheMock.AssertCalled(t, "Get", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything)
	cacheMock.AssertCalled(t, "Set", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything, 5*time.Minute)
}

func TestShouldNotPutValueToCacheWhenEmptyIfResponseError(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Get", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything).Return(nil, nil)
	cacheMock.On("Set", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything, 5*time.Minute).Return(nil)

	mw := Cacheable(cacheMock, "addresses")(func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "cached empty"), errors.New("response error")
	})

	resp, err := mw(nil, "params")

	assert.EqualError(t, err, "response error")
	assert.Equal(t, "cached empty", resp.Data())
	cacheMock.AssertCalled(t, "Get", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything)
	cacheMock.AssertNotCalled(t, "Set", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything, 0)
}

func TestShouldNotPutValueToCacheWhenEmptyIfResponseNil(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Get", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything).Return(nil, nil)
	cacheMock.On("Set", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything, 5*time.Minute).Return(nil)

	mw := Cacheable(cacheMock, "addresses", CacheOptions{TTL: 5 * time.Minute})(func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return nil, errors.New("response error")
	})

	resp, err := mw(nil, "params")

	assert.EqualError(t, err, "response error")
	assert.Nil(t, resp)
	cacheMock.AssertCalled(t, "Get", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything)
	cacheMock.AssertNotCalled(t, "Set", "0b60ebc36c1f4764fde1db73790133fe", mock.Anything, 5*time.Minute)
}
