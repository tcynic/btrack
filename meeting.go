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

	meetingID, err := a.meetings.Create(
		input.ProjectID,
		input.Title,
		input.MeetingDate,
		input.DurationMinutes,
		input.Attendees,
		input.Notes,
	)
	if err != nil {
		return nil, err
	}

	return a.GetMeeting(meetingID)
}

// GetMeetings returns all meetings for a project
func (a *App) GetMeetings(projectID int64) ([]models.Meeting, error) {
	meetings, err := a.meetings.GetByProject(projectID)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(meetings), nil
}

// GetMeeting returns a single meeting by ID
func (a *App) GetMeeting(id int64) (*models.Meeting, error) {
	return a.meetings.GetByID(id)
}

// UpdateMeeting updates an existing meeting
func (a *App) UpdateMeeting(input models.UpdateMeetingInput) (*models.Meeting, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	err := a.meetings.Update(
		input.ID,
		input.Title,
		input.MeetingDate,
		input.DurationMinutes,
		input.Attendees,
		input.Notes,
	)
	if err != nil {
		return nil, err
	}

	return a.GetMeeting(input.ID)
}

// DeleteMeeting removes a meeting
func (a *App) DeleteMeeting(id int64) error {
	// First delete associated tasks
	if err := a.tasks.DeleteBySource("meeting", id); err != nil {
		return err
	}

	return a.meetings.Delete(id)
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
	searchPattern := "%" + query + "%"
	meetings, err := a.meetings.Search(searchPattern)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(meetings), nil
}
