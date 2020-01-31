package endpoint

// EndpointResponse the request result.
type EndpointResponse interface {
	Code() int
	Data() interface{}
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

func (e *endpointResponse) Code() int {
	return e.code
}

func (e *endpointResponse) Data() interface{} {
	return e.Data()
}
