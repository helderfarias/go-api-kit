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

func TestDeleteAllEntriesWhenCacheEvict(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("DeleteAll", "addresses").Return(nil)

	service := func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "return after, not cached"), nil
	}

	mw := CacheEvict(cacheMock, "addresses", CacheEvictOptions{AllEntries: true})(service)

	resp, err := mw(nil, "request")

	assert.Nil(t, err)
	assert.Equal(t, "return after, not cached", resp.Data())
	cacheMock.AssertExpectations(t)
}

func TestDeleteOnlyEntryWhenCacheEvict(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Delete", "addresses").Return(nil)

	service := func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "return after, not cached"), nil
	}

	mw := CacheEvict(cacheMock, "addresses")(service)

	resp, err := mw(nil, "request")

	assert.Nil(t, err)
	assert.Equal(t, "return after, not cached", resp.Data())
	cacheMock.AssertExpectations(t)
}

func TestDeleteAndSetCacheWhenCachePut(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Delete", "addresses").Return(nil)
	cacheMock.On("Set", "addresses:bc2b9fed1ade259444436fb721d3ab22", mock.Anything, mock.Anything).Return(nil)

	service := func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "return after cached"), nil
	}

	mw := CachePut(cacheMock, "addresses")(service)

	resp, err := mw(nil, "request")

	assert.Nil(t, err)
	assert.Equal(t, "return after cached", resp.Data())
	cacheMock.AssertExpectations(t)
}

func TestDeleteCacheWhenCachePutIfServiceResultError(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Delete", "addresses").Return(nil)

	service := func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, nil), errors.New("service erro")
	}

	mw := CachePut(cacheMock, "addresses")(service)

	resp, err := mw(nil, "request")

	assert.EqualError(t, err, "service erro")
	assert.Nil(t, resp.Data())
	cacheMock.AssertNotCalled(t, "Set", mock.Anything, mock.Anything)
	cacheMock.AssertExpectations(t)
}

func TestValueFromCacheWhenNotEmpty(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Get", "addresses:bc2b9fed1ade259444436fb721d3ab22", mock.Anything).Return(&entryCache{200, "cached: address 10"}, nil)

	mw := Cacheable(cacheMock, "addresses")(func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "not return"), nil
	})

	resp, err := mw(nil, "request")

	assert.Nil(t, err)
	assert.Equal(t, "cached: address 10", resp.Data())
	cacheMock.AssertExpectations(t)
}

func TestPutValueToCacheWhenEmpty(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Get", "addresses:8a80b0b2fc5b41f18697478e5031ca22", mock.Anything).Return(nil, nil)
	cacheMock.On("Set", "addresses:8a80b0b2fc5b41f18697478e5031ca22", mock.Anything, 5*time.Minute).Return(nil)

	mw := Cacheable(cacheMock, "addresses", CacheableOptions{TTL: 5 * time.Minute})(func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "cached empty"), nil
	})

	resp, err := mw(nil, "params")

	assert.Nil(t, err)
	assert.Equal(t, "cached empty", resp.Data())
	cacheMock.AssertExpectations(t)
}

func TestShouldNotPutValueToCacheWhenEmptyIfResponseError(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Get", "addresses:8a80b0b2fc5b41f18697478e5031ca22", mock.Anything).Return(nil, nil)

	mw := Cacheable(cacheMock, "addresses")(func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return endpoint.Response(200, "cached empty"), errors.New("response error")
	})

	resp, err := mw(nil, "params")

	assert.EqualError(t, err, "response error")
	assert.Equal(t, "cached empty", resp.Data())
	cacheMock.AssertExpectations(t)
}

func TestShouldNotPutValueToCacheWhenEmptyIfResponseNil(t *testing.T) {
	cacheMock := &cacheServerMock{}

	cacheMock.On("Get", "addresses:8a80b0b2fc5b41f18697478e5031ca22", mock.Anything).Return(nil, nil)

	mw := Cacheable(cacheMock, "addresses", CacheableOptions{TTL: 5 * time.Minute})(func(ctx context.Context, request interface{}) (endpoint.EndpointResponse, error) {
		return nil, errors.New("response error")
	})

	resp, err := mw(nil, "params")

	assert.EqualError(t, err, "response error")
	assert.Nil(t, resp)
	cacheMock.AssertExpectations(t)
}
