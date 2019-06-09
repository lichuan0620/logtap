package httputil

import "net/http"

type httpError struct {
	Code    int    `json:"code"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// NewHTTPError returns an error that, when unmarshalled, shows enriched HTTP error information.
func NewHTTPError(code int, reason, message string) error {
	return &httpError{
		Code:    code,
		Reason:  reason,
		Message: message,
	}
}

// Error implements the error interface and returns an abbreviated error message.
func (e httpError) Error() string {
	return e.Reason
}

// NewNotFoundError should be used when the requested resource cannot be found.
func NewNotFoundError(message string) error {
	const reason = "target not found"
	return NewHTTPError(http.StatusNotFound, reason, message)
}

// NewMethodNotAllowedError should be used when the method used is not supported by the requested resource.
func NewMethodNotAllowedError() error {
	const reason = "method not allowed"
	return NewHTTPError(http.StatusMethodNotAllowed, reason, "")
}

// NewRequestError should be used when the client made a invalid request.
func NewRequestError(message string) error {
	const reason = "invalid request"
	return NewHTTPError(http.StatusBadRequest, reason, message)
}

// NewValidationError should be used when a requested resource is found but not valid.
func NewValidationError(message string) error {
	const reason = "invalid data"
	return NewHTTPError(http.StatusInternalServerError, reason, message)
}

// NewNotImplementedError should be used when a request requires a feature that have not been implemented
func NewNotImplementedError() error {
	const reason = "not implemented"
	return NewHTTPError(http.StatusNotImplemented, reason, "")
}
