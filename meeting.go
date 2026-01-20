package main

import (
	"fmt"

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
		m, err := models.ScanMeeting(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan meeting: %w", err)
		}
		meetings = append(meetings, *m)
	}

	return database.EnsureSlice(meetings), nil
}

// GetMeeting returns a single meeting by ID
func (a *App) GetMeeting(id int64) (*models.Meeting, error) {
	row := a.db.QueryRow(database.SelectMeetingByID, id)
	m, err := models.ScanMeeting(row.Scan)
	if err != nil {
		return nil, models.NotFound("meeting")
	}
	return m, nil
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
		return nil, models.NotFound("meeting")
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
		return models.NotFound("meeting")
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
		m, err := models.ScanMeetingWithProject(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan meeting: %w", err)
		}
		meetings = append(meetings, *m)
	}

	return database.EnsureSlice(meetings), nil
}

// GetMeetingsByWeek returns all meetings for a week (7-day period starting from weekStartDate)
func (a *App) GetMeetingsByWeek(weekStartDate string) ([]models.MeetingWithProject, error) {
	// Calculate the end date (7 days after start)
	startDate, err := database.ParseDate(weekStartDate)
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
		m, err := models.ScanMeetingWithProject(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan meeting: %w", err)
		}
		meetings = append(meetings, *m)
	}

	return database.EnsureSlice(meetings), nil
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
		m, err := models.ScanMeetingWithProject(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan meeting: %w", err)
		}
		meetings = append(meetings, *m)
	}

	return database.EnsureSlice(meetings), nil
}
