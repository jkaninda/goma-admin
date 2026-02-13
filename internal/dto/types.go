package dto

type ResponseDto struct {
	Success    bool   `json:"success"`
	StatusCode string `json:"statusCode"`
	Details    any    `json:"details,omitempty"`
}

type Response[T any] struct {
	ResponseDto
	Data T `json:"data,omitempty"`
}

func SuccessResponse[T any](message string, data T) Response[T] {
	return Response[T]{
		ResponseDto: ResponseDto{
			Success: true,
		},
		Data: data,
	}
}
func ErrorResponse(message string, err error) Response[any] {
	return Response[any]{
		ResponseDto: ResponseDto{
			Success: false,
			Details: err.Error(),
		},
	}
}
func ErrorResponseData(message string, err error, data any) Response[any] {
	return Response[any]{
		ResponseDto: ResponseDto{
			Success: false,
			Details: err,
		},
		Data: data,
	}
}
