package dto

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

func Success[T any](data T) SuccessResponse[T] {
	return SuccessResponse[T]{Data: data}
}
