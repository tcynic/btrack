package models

// Meeting represents a project meeting
type Meeting struct {
	ID              int64  `json:"id"`
	ProjectID       int64  `json:"projectId"`
	Title           string `json:"title"`
	MeetingDate     string `json:"meetingDate"`
	DurationMinutes int    `json:"durationMinutes"`
	Attendees       string `json:"attendees"`
	Notes           string `json:"notes"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// ScanMeeting scans a database row into a Meeting struct.
// Expected columns: id, project_id, title, meeting_date, duration_minutes, attendees, notes, created_at, updated_at
func ScanMeeting(scan func(dest ...any) error) (*Meeting, error) {
	var m Meeting
	var attendees, notes *string
	var createdAt, updatedAt string

	err := scan(
		&m.ID,
		&m.ProjectID,
		&m.Title,
		&m.MeetingDate,
		&m.DurationMinutes,
		&attendees,
		&notes,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	m.Attendees = NullableString(attendees)
	m.Notes = NullableString(notes)
	m.CreatedAt = createdAt
	m.UpdatedAt = updatedAt

	return &m, nil
}

// ScanMeetingWithProject scans a database row into a MeetingWithProject struct.
// Expected columns: id, project_id, title, meeting_date, duration_minutes, attendees, notes, created_at, updated_at, project_name
func ScanMeetingWithProject(scan func(dest ...any) error) (*MeetingWithProject, error) {
	var m MeetingWithProject
	var attendees, notes *string
	var createdAt, updatedAt string

	err := scan(
		&m.ID,
		&m.ProjectID,
		&m.Title,
		&m.MeetingDate,
		&m.DurationMinutes,
		&attendees,
		&notes,
		&createdAt,
		&updatedAt,
		&m.ProjectName,
	)
	if err != nil {
		return nil, err
	}

	m.Attendees = NullableString(attendees)
	m.Notes = NullableString(notes)
	m.CreatedAt = createdAt
	m.UpdatedAt = updatedAt

	return &m, nil
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
		return ValidationError("title", "title is required")
	}
	if c.MeetingDate == "" {
		return ValidationError("meetingDate", "date is required")
	}
	if c.DurationMinutes <= 0 {
		c.DurationMinutes = 60 // Default to 60 minutes
	}
	return nil
}

// Validate checks if the UpdateMeetingInput is valid
func (u *UpdateMeetingInput) Validate() error {
	if u.Title == "" {
		return ValidationError("title", "title is required")
	}
	if u.MeetingDate == "" {
		return ValidationError("meetingDate", "date is required")
	}
	if u.DurationMinutes <= 0 {
		u.DurationMinutes = 60
	}
	return nil
}
