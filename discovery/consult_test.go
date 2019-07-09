package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRegister(t *testing.T) {
	options := &ConsulClientSettigs{}

	register := NewConsulRegister(options)

	assert.NotNil(t, register)
}
