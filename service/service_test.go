package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMiddleware(t *testing.T) {
	s := func(next Service) Service {
		return func(ctx context.Context) (interface{}, error) {
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
