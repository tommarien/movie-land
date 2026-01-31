package api

import (
	"fmt"
	"net/http"
)

type HttpStatusError struct {
	statusCode int
	message    string
	errors     []string
}

func NewNotFoundStatusError(message string) *HttpStatusError {
	if message == "" {
		message = "resource not found"
	}

	return &HttpStatusError{
		statusCode: http.StatusNotFound,
		message:    message,
	}
}

func NewBadRequestStatusError(message string, errors []string) *HttpStatusError {
	if message == "" {
		message = "bad request"
	}

	return &HttpStatusError{
		statusCode: http.StatusBadRequest,
		message:    message,
		errors:     errors,
	}
}

func NewConflictStatusError(message string) *HttpStatusError {
	if message == "" {
		message = "conflict"
	}

	return &HttpStatusError{
		statusCode: http.StatusConflict,
		message:    message,
	}
}

func (e *HttpStatusError) Error() string {
	return fmt.Sprintf("http %d: %s", e.statusCode, e.message)
}

func (e *HttpStatusError) StatusCode() int {
	return e.statusCode
}

func (e *HttpStatusError) Message() string {
	return e.message
}

func (e *HttpStatusError) Errors() []string {
	return e.errors
}
