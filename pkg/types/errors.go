// Package types provides shared types for API responses and errors.
//
// Error strategy: APIError and the predefined Err* values (ErrNotFound, ErrUnauthorized, etc.)
// are used when mapping domain errors to HTTP responses. Controllers map service/domain errors
// (e.g. auth.Err*) to types.APIError and send the response via utils.RespondWithAPIError.
// All HTTP error responses go through pkg/utils (BadRequest, Unauthorized, RespondWithAPIError, etc.);
// do not duplicate status-code logic outside the controller mapping and utils.
package types

import (
	"fmt"
	"net/http"
)

// APIError represents a standardized API error response
type APIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// NewAPIError creates a new APIError
func NewAPIError(code int, message string, details interface{}) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Predefined common errors
var (
	ErrInvalidInput      = NewAPIError(http.StatusBadRequest, "Invalid input data", nil)
	ErrNotFound          = NewAPIError(http.StatusNotFound, "Resource not found", nil)
	ErrUnauthorized      = NewAPIError(http.StatusUnauthorized, "Unauthorized access", nil)
	ErrForbidden         = NewAPIError(http.StatusForbidden, "Access forbidden", nil)
	ErrConflict          = NewAPIError(http.StatusConflict, "Resource already exists", nil)
	ErrInternalServer    = NewAPIError(http.StatusInternalServerError, "Internal server error", nil)
	ErrDatabaseError     = NewAPIError(http.StatusInternalServerError, "Database operation failed", nil)
	ErrExternalService   = NewAPIError(http.StatusBadGateway, "External service error", nil)
	ErrRateLimitExceeded = NewAPIError(http.StatusTooManyRequests, "Rate limit exceeded", nil)
)

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Error implements the error interface
func (ve ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", ve.Message)
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	if len(ve) == 1 {
		return fmt.Sprintf("validation error: %s", ve[0].Message)
	}
	return fmt.Sprintf("validation errors: %d errors found", len(ve))
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewValidationErrorWithValue creates a new validation error with value
func NewValidationErrorWithValue(field, message, value string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	}
}
