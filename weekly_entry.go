package main

import (
	"btrack/internal/models"
)

// GetWeeklyEntries returns all weekly entries for a project with status information
func (a *App) GetWeeklyEntries(projectID int64) ([]models.WeeklyEntryWithStatus, error) {
	return a.trackingService.GetWeeklyEntries(projectID)
}

// UpdateActualHours updates the actual hours for a specific weekly entry
func (a *App) UpdateActualHours(input models.UpdateActualHoursInput) (*models.WeeklyEntryWithStatus, error) {
	return a.trackingService.UpdateActualHours(input)
}

// GetWeeklyEntriesByWeek returns all weekly entries for a specific week across all active projects
func (a *App) GetWeeklyEntriesByWeek(weekStartDate string) ([]models.WeeklyEntryWithProject, error) {
	return a.trackingService.GetWeeklyEntriesByWeek(weekStartDate)
}
