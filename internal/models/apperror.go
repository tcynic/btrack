package models

import "fmt"

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	ErrCodeNotFound    ErrorCode = "NOT_FOUND"
	ErrCodeValidation  ErrorCode = "VALIDATION"
	ErrCodeDatabase    ErrorCode = "DATABASE"
	ErrCodeConflict    ErrorCode = "CONFLICT"
	ErrCodeForbidden   ErrorCode = "FORBIDDEN"
	ErrCodeInternal    ErrorCode = "INTERNAL"
)

// AppError represents a structured application error
type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"`
	Err     error     `json:"-"` // Internal error, not exposed to clients
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NotFound creates a NOT_FOUND error
func NotFound(entity string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", entity),
	}
}

// ValidationError creates a VALIDATION error
func ValidationError(field, message string) *AppError {
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
		Field:   field,
	}
}

// DatabaseError creates a DATABASE error
func DatabaseError(err error) *AppError {
	return &AppError{
		Code:    ErrCodeDatabase,
		Message: "Database operation failed",
		Err:     err,
	}
}

// Conflict creates a CONFLICT error
func Conflict(message string) *AppError {
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
	}
}

// Forbidden creates a FORBIDDEN error
func Forbidden(message string) *AppError {
	return &AppError{
		Code:    ErrCodeForbidden,
		Message: message,
	}
}

// Internal creates an INTERNAL error
func Internal(err error, message string) *AppError {
	return &AppError{
		Code:    ErrCodeInternal,
		Message: message,
		Err:     err,
	}
}
