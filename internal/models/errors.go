package models

import "errors"

// NullableString converts a nullable string pointer to a string
// Returns empty string if the pointer is nil
func NullableString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Error variables for validation
var (
	ErrInvalidStatus     = errors.New("invalid status value")
	ErrInvalidPriority   = errors.New("invalid priority value")
	ErrInvalidSourceType = errors.New("invalid source type")
)
