package mb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNatsServerWithSubAndPubEmpty(t *testing.T) {
	ns := NewNatsServer()

	assert.IsType(t, &emptyPub{}, ns.Pub())
	assert.IsType(t, &emptySub{}, ns.Sub())
}
