package validate

import "errors"

// ErrInvalidID occurs when an ID is not in a valid format.
var ErrInvalidID = errors.New("ID is not in its proper form")

// ErrorResponse is the format used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string `json:"error"`
	Fields string `json:"fileds,omitempty"`
}

/**
Request Error is used to pass an error during the request through the application
with web specific context
*/
type RequestError struct {
	Err    error
	Status int
	Fields error
}

/**
Error implements the error interface. It uses the default message of the wrapped error. This is what will
be shown in the services logs.
*/
func (err *RequestError) Error() string {
	return err.Err.Error()
}

/**
NewRequestError wraps a provided error with an HTTP status code. This function should be used when handlers
encounter expected errors.
*/
func NewRequestError(err error, status int) error {
	return &RequestError{err, status, nil}
}
