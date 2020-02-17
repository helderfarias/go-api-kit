package endpoint

// EndpointResponse the request result.
type EndpointResponse interface {
	Code() int
	Data() interface{}
}

type paging struct {
	Page  int64 `json:"page"`
	Total int64 `json:"total"`
	Limit int64 `json:"limit"`
}

type entityPaging struct {
	Data   interface{} `json:"data"`
	Paging paging      `json:"paging"`
}

type endpointResponse struct {
	code int
	data interface{}
}

// Response transfer object
func Response(code int, data interface{}) EndpointResponse {
	return &endpointResponse{
		code: code,
		data: data,
	}
}

// Paginate transfer object
func Paginate(data interface{}, page, limit, total int64) interface{} {
	paging := paging{}

	if total != 0 {
		paging.Total = total
	}

	if limit != 0 {
		paging.Limit = limit
	}

	if page != 0 {
		paging.Page = page
	}

	return entityPaging{Data: data, Paging: paging}
}

func (e *endpointResponse) Code() int {
	return e.code
}

func (e *endpointResponse) Data() interface{} {
	return e.data
}
