package middleware

import (
	"context"
	"testing"

	"github.com/helderfarias/go-api-kit/endpoint"
	"github.com/stretchr/testify/assert"
)

func TestCreateMiddleware(t *testing.T) {
	s := func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(parent context.Context, request interface{}) (response endpoint.EndpointResponse, err error) {
			return nil, nil
		}
	}

	assert.NotNil(t, s)
}

func TestCreateDatatabaseMiddleware(t *testing.T) {
	assert.NotNil(t, Database(nil, ""))
}

func TestCreateDatatabaseTxMiddleware(t *testing.T) {
	assert.NotNil(t, DatabaseWithTx(nil, ""))
}
