package apperr

import (
	"errors"
	"fmt"
	"net/http"
)

type Error struct {
	Code       string
	Message    string
	HTTPStatus int
	Details    map[string]any
	Err        error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return e.Code
}

func (e *Error) Unwrap() error { return e.Err }

func New(code, message string, status int) *Error {
	return &Error{Code: code, Message: message, HTTPStatus: status}
}

func Validation(message string, details map[string]any) *Error {
	return &Error{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Details:    details,
	}
}

func From(err error) *Error {
	var ae *Error
	if errors.As(err, &ae) {
		return ae
	}
	return &Error{
		Code:       "INTERNAL_ERROR",
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

var (
	ErrItemNotFound   = New("ITEM_NOT_FOUND", "Item not found", http.StatusNotFound)
	ErrNotImplemented = New("NOT_IMPLEMENTED", "Not implemented in this phase", http.StatusNotImplemented)
)
