package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresConnectionFactory(t *testing.T) {
	f, err := NewPostgresConnectionFactory("user=thello password=thello00 dbname=thello host=localhost port=5432 sslmode=disable", 1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, f)
}
