// internal/domain/models/errors.go
package models

import "net/http"

// ErrorType is a custom type to represent different kinds of errors.
type ErrorType string

// Enumeration of custom error types for consistent error handling across the application
const (
	ValidationError     ErrorType = "VALIDATION_ERROR"
	ServerError         ErrorType = "SERVER_ERROR"
	AuthenticationError ErrorType = "AUTHENTICATION_ERROR"
	NotFoundError       ErrorType = "NOT_FOUND_ERROR"
	MethodNotAllowed    ErrorType = "METHOD_NOT_ALLOWED"
	ForbiddenError      ErrorType = "FORBIDDEN_ERROR"
	UnauthorizedError   ErrorType = "UNAUTHORIZED_ERROR"
	TimeoutError        ErrorType = "TIMEOUT"
	BadRequestError     ErrorType = "BAD_REQUEST_ERROR"
	ProxyError          ErrorType = "PROXY_ERROR" // Added ProxyError
	LoggingError        ErrorType = "LOGGING_ERROR"
)

// statusCodeErrorMap maps status codes to corresponding error types, status messages, and default messages
var statusCodeErrorMap = map[int]struct {
	ErrorType     ErrorType
	StatusMessage string
	Message       string
}{
	http.StatusBadRequest:          {BadRequestError, "Bad Request", "The request was invalid or cannot be served"},
	http.StatusUnauthorized:        {UnauthorizedError, "Unauthorized", "The request requires user authentication"},
	http.StatusForbidden:           {ForbiddenError, "Forbidden", "Access is forbidden to the requested resource"},
	http.StatusNotFound:            {NotFoundError, "Not Found", "The requested resource could not be found"},
	http.StatusMethodNotAllowed:    {MethodNotAllowed, "Method Not Allowed", "The requested method is not allowed for the resource"},
	http.StatusRequestTimeout:      {TimeoutError, "Request Timeout", "The server timed out waiting for the request"},
	http.StatusInternalServerError: {ServerError, "Internal Server Error", "The server encountered an internal error"},
	http.StatusGatewayTimeout:      {ProxyError, "Gateway Timeout", "The server, while acting as a gateway, timed out"},
	// Add more status codes and messages as needed
}

func StatusCodeToApplicationError(statusCode int) *ApplicationError {
	if errInfo, ok := statusCodeErrorMap[statusCode]; ok {
		return NewApplicationError(
			errInfo.ErrorType,
			errInfo.Message,
			statusCode,
			nil, // Passing nil to errors, as message will be part of "message" field
		)
	}

	// Default to ServerError if status code is unhandled
	return NewApplicationError(
		ServerError,
		"An unknown server error occurred",
		statusCode,
		nil,
	)
}