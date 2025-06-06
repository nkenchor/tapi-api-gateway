package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ApplicationError represents a structured error for the application
type ApplicationError struct {
	ErrorType      ErrorType           `json:"error_type"`
	ErrorReference string              `json:"error_reference"`
	Timestamp      string              `json:"timestamp"`
	StatusCode     int                 `json:"status_code"`
	Errors         map[string][]string `json:"errors"`
}

// NewApplicationError creates a new instance of ApplicationError
func NewApplicationError(errorType ErrorType, message string, statusCode int, errors map[string][]string) *ApplicationError {
	if errors == nil {
		errors = map[string][]string{
			"message": {message},
		}
	}
	return &ApplicationError{
		ErrorType:      errorType,
		ErrorReference: uuid.New().String(),
		Timestamp:      time.Now().Format(time.RFC3339),
		StatusCode:     statusCode,
		Errors:         errors,
	}
}

// ToJSON converts the ApplicationError to a JSON string
func (e *ApplicationError) ToJSON() string {
	errorJSON, _ := json.Marshal(e)
	return string(errorJSON)
}

// Error implements the error interface
func (e *ApplicationError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Timestamp, e.ErrorType, e.Errors["message"][0])
}
