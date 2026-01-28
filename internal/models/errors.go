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

// Validation errors
var (
	ErrNameRequired           = errors.New("project name is required")
	ErrInvalidHours           = errors.New("hours must be a positive number")
	ErrSpecialistHoursTooHigh = errors.New("specialist hours must be less than total sold hours")
	ErrDatesRequired          = errors.New("start date and end date are required")
	ErrInvalidDateRange       = errors.New("start date must be before end date")
	ErrProjectNotFound        = errors.New("project not found")
	ErrEntryNotFound          = errors.New("weekly entry not found")
	ErrCannotEditFutureWeek   = errors.New("cannot edit actual hours for future weeks")

	// Meeting, Note, Goal errors
	ErrTitleRequired   = errors.New("title is required")
	ErrDateRequired    = errors.New("date is required")
	ErrMeetingNotFound = errors.New("meeting not found")
	ErrNoteNotFound    = errors.New("note not found")
	ErrGoalNotFound    = errors.New("goal not found")
	ErrInvalidStatus   = errors.New("invalid status value")

	// Task errors
	ErrTaskNotFound        = errors.New("task not found")
	ErrInvalidPriority     = errors.New("invalid priority value")
	ErrInvalidSourceType   = errors.New("invalid source type")
)
