package models

import "time"

// Meeting represents a project meeting
type Meeting struct {
	ID              int64     `json:"id"`
	ProjectID       int64     `json:"projectId"`
	Title           string    `json:"title"`
	MeetingDate     string    `json:"meetingDate"`
	DurationMinutes int       `json:"durationMinutes"`
	Attendees       string    `json:"attendees"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// MeetingWithProject extends Meeting with project name for cross-project queries
type MeetingWithProject struct {
	Meeting
	ProjectName string `json:"projectName"`
}

// CreateMeetingInput is the input for creating a new meeting
type CreateMeetingInput struct {
	ProjectID       int64  `json:"projectId"`
	Title           string `json:"title"`
	MeetingDate     string `json:"meetingDate"`
	DurationMinutes int    `json:"durationMinutes"`
	Attendees       string `json:"attendees"`
	Notes           string `json:"notes"`
}

// UpdateMeetingInput is the input for updating a meeting
type UpdateMeetingInput struct {
	ID              int64  `json:"id"`
	Title           string `json:"title"`
	MeetingDate     string `json:"meetingDate"`
	DurationMinutes int    `json:"durationMinutes"`
	Attendees       string `json:"attendees"`
	Notes           string `json:"notes"`
}

// Validate checks if the CreateMeetingInput is valid
func (c *CreateMeetingInput) Validate() error {
	if c.Title == "" {
		return ErrTitleRequired
	}
	if c.MeetingDate == "" {
		return ErrDateRequired
	}
	if c.DurationMinutes <= 0 {
		c.DurationMinutes = 60 // Default to 60 minutes
	}
	return nil
}

// Validate checks if the UpdateMeetingInput is valid
func (u *UpdateMeetingInput) Validate() error {
	if u.Title == "" {
		return ErrTitleRequired
	}
	if u.MeetingDate == "" {
		return ErrDateRequired
	}
	if u.DurationMinutes <= 0 {
		u.DurationMinutes = 60
	}
	return nil
}
