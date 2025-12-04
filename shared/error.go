// shared/errors.go
package shared

import "net/http"

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewBadRequest(msg string) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    msg,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewNotFound(msg string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    msg,
		HTTPStatus: http.StatusNotFound,
	}
}

func NewInternal(msg string) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    msg,
		HTTPStatus: http.StatusInternalServerError,
	}
}
