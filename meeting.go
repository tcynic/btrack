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
	// First delete associated tasks
	_, err := a.db.Exec(database.DeleteTasksBySource, "meeting", id)
	if err != nil {
		return fmt.Errorf("failed to delete associated tasks: %w", err)
	}

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

// GetMeetingsByDate returns all meetings for a specific date across all projects
func (a *App) GetMeetingsByDate(date string) ([]models.MeetingWithProject, error) {
	rows, err := a.db.Query(database.SelectMeetingsByDate, date)
	if err != nil {
		return nil, fmt.Errorf("failed to query meetings by date: %w", err)
	}
	defer rows.Close()

	var meetings []models.MeetingWithProject
	for rows.Next() {
		var m models.MeetingWithProject
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
			&m.ProjectName,
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
		meetings = []models.MeetingWithProject{}
	}

	return meetings, nil
}

// GetMeetingsByWeek returns all meetings for a week (7-day period starting from weekStartDate)
func (a *App) GetMeetingsByWeek(weekStartDate string) ([]models.MeetingWithProject, error) {
	// Calculate the end date (7 days after start)
	startDate, err := time.Parse("2006-01-02", weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid week start date: %w", err)
	}
	endDate := startDate.AddDate(0, 0, 7)

	rows, err := a.db.Query(database.SelectMeetingsByWeek, weekStartDate, endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query meetings by week: %w", err)
	}
	defer rows.Close()

	var meetings []models.MeetingWithProject
	for rows.Next() {
		var m models.MeetingWithProject
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
			&m.ProjectName,
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
		meetings = []models.MeetingWithProject{}
	}

	return meetings, nil
}

// SearchMeetings searches for meetings by title, notes, or attendees
func (a *App) SearchMeetings(query string) ([]models.MeetingWithProject, error) {
	rows, err := a.db.Query(database.SearchMeetings, query, query, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search meetings: %w", err)
	}
	defer rows.Close()

	var meetings []models.MeetingWithProject
	for rows.Next() {
		var m models.MeetingWithProject
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
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.ProjectName,
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
		meetings = []models.MeetingWithProject{}
	}

	return meetings, nil
}
