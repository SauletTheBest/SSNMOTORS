package dto

import "net/http"

type HTTPError struct {
	Code    int
	Message string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func FromError(err error) *HTTPError {
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}
