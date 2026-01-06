package models

import "errors"

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
)
