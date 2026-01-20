package database

import (
	"fmt"
	"time"
)

// ParseDate parses a date string (YYYY-MM-DD) into a time.Time.
// Returns an error if parsing fails.
func ParseDate(s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}
	return t, nil
}

// NullableString returns the dereferenced string value or empty string if nil.
func NullableString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// NullableInt64 returns the dereferenced int64 value or nil pointer if nil.
func NullableInt64(i *int64) *int64 {
	return i
}

// EnsureSlice returns an empty slice of the same type if the input is nil.
// This prevents JSON null responses when there are no results.
func EnsureSlice[T any](slice []T) []T {
	if slice == nil {
		return []T{}
	}
	return slice
}
