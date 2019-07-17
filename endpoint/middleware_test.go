package endpoint

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMiddleware(t *testing.T) {
	s := func(next Endpoint) Endpoint {
		return func(parent context.Context, request interface{}) (response interface{}, err error) {
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
