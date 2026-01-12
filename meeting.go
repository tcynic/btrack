package main

import (
	"fmt"
	"time"

	"btrack/internal/database"
	"btrack/internal/models"
)

// CreateMeeting creates a new meeting for a project
func (a *App) CreateMeeting(input models.CreateMeetingInput) (*models.Meeting, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	result, err := a.db.Exec(database.InsertMeeting,
		input.ProjectID,
		input.Title,
		input.MeetingDate,
		input.DurationMinutes,
		input.Attendees,
		input.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert meeting: %w", err)
	}

	meetingID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get meeting ID: %w", err)
	}

	return a.GetMeeting(meetingID)
}

// GetMeetings returns all meetings for a project
func (a *App) GetMeetings(projectID int64) ([]models.Meeting, error) {
	rows, err := a.db.Query(database.SelectMeetingsByProject, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query meetings: %w", err)
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var m models.Meeting
		var attendees, notes *string
		var createdAt, updatedAt string

		err := rows.Scan(
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
			return nil, fmt.Errorf("failed to scan meeting: %w", err)
		}

		if attendees != nil {
			m.Attendees = *attendees
		}
		if notes != nil {
			m.Notes = *notes
		}
		m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		m.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		meetings = append(meetings, m)
	}

	if meetings == nil {
		meetings = []models.Meeting{}
	}

	return meetings, nil
}

// GetMeeting returns a single meeting by ID
func (a *App) GetMeeting(id int64) (*models.Meeting, error) {
	var m models.Meeting
	var attendees, notes *string
	var createdAt, updatedAt string

	err := a.db.QueryRow(database.SelectMeetingByID, id).Scan(
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
		return nil, models.ErrMeetingNotFound
	}

	if attendees != nil {
		m.Attendees = *attendees
	}
	if notes != nil {
		m.Notes = *notes
	}
	m.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	m.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &m, nil
}

// UpdateMeeting updates an existing meeting
func (a *App) UpdateMeeting(input models.UpdateMeetingInput) (*models.Meeting, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	result, err := a.db.Exec(database.UpdateMeeting,
		input.Title,
		input.MeetingDate,
		input.DurationMinutes,
		input.Attendees,
		input.Notes,
		input.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update meeting: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, models.ErrMeetingNotFound
	}

	return a.GetMeeting(input.ID)
}

// DeleteMeeting removes a meeting
func (a *App) DeleteMeeting(id int64) error {
	result, err := a.db.Exec(database.DeleteMeeting, id)
	if err != nil {
		return fmt.Errorf("failed to delete meeting: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrMeetingNotFound
	}

	return nil
}
