package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an error that carries an HTTP status code and a safe
// user-facing message.
type AppError struct {
	StatusCode int
	Message    string
	Err        error
}

func NewAppError(statusCode int, message string) *AppError {
	return &AppError{StatusCode: statusCode, Message: message}
}

func NewAppErrorf(statusCode int, format string, args ...interface{}) *AppError {
	return &AppError{StatusCode: statusCode, Message: fmt.Sprintf(format, args...)}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

// IsAppError reports whether the error is an AppError instance.
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok && appErr != nil
}

// HTTPStatus maps an error to its HTTP status code. Unknown errors map to 500.
func HTTPStatus(err error) int {
	if appErr, ok := IsAppError(err); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}

// Message returns the safe user-facing message of an error.
func Message(err error) string {
	if appErr, ok := IsAppError(err); ok {
		return appErr.Message
	}
	return http.StatusText(HTTPStatus(err))
}
