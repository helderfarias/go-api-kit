package endpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type dataFake struct {
}

func TestResponseWithPagination(t *testing.T) {
	data := &dataFake{}

	resp := Response(200, Paginate(data, 1, 10, 20))

	assert.Equal(t, entityPaging{
		Data: data,
		Paging: paging{
			Page:  1,
			Total: 20,
			Limit: 10,
		},
	}, resp.Data())
}
