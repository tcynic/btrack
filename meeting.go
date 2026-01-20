package main

import (
	"btrack/internal/models"
)

// CreateMeeting creates a new meeting for a project
func (a *App) CreateMeeting(input models.CreateMeetingInput) (*models.Meeting, error) {
	return a.notesService.CreateMeeting(input)
}

// GetMeetings returns all meetings for a project
func (a *App) GetMeetings(projectID int64) ([]models.Meeting, error) {
	return a.notesService.GetMeetings(projectID)
}

// GetMeeting returns a single meeting by ID
func (a *App) GetMeeting(id int64) (*models.Meeting, error) {
	return a.notesService.GetMeeting(id)
}

// UpdateMeeting updates an existing meeting
func (a *App) UpdateMeeting(input models.UpdateMeetingInput) (*models.Meeting, error) {
	return a.notesService.UpdateMeeting(input)
}

// DeleteMeeting removes a meeting
func (a *App) DeleteMeeting(id int64) error {
	return a.notesService.DeleteMeeting(id)
}

// GetMeetingsByDate returns all meetings for a specific date across all projects
func (a *App) GetMeetingsByDate(date string) ([]models.MeetingWithProject, error) {
	return a.notesService.GetMeetingsByDate(date)
}

// GetMeetingsByWeek returns all meetings for a week (7-day period starting from weekStartDate)
func (a *App) GetMeetingsByWeek(weekStartDate string) ([]models.MeetingWithProject, error) {
	return a.notesService.GetMeetingsByWeek(weekStartDate)
}

// SearchMeetings searches for meetings by title, notes, or attendees
func (a *App) SearchMeetings(query string) ([]models.MeetingWithProject, error) {
	return a.notesService.SearchMeetings(query)
}
