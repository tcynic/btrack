package main

import (
	"fmt"
	"time"

	"btrack/internal/database"
	"btrack/internal/models"
)

// GetWeeklyEntries returns all weekly entries for a project with status information
func (a *App) GetWeeklyEntries(projectID int64) ([]models.WeeklyEntryWithStatus, error) {
	rows, err := a.db.Query(database.SelectWeeklyEntriesByProject, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query weekly entries: %w", err)
	}
	defer rows.Close()

	currentMonday := getCurrentWeekMonday()
	var entries []models.WeeklyEntryWithStatus

	for rows.Next() {
		var e models.WeeklyEntry
		var createdAt, updatedAt string

		err := rows.Scan(
			&e.ID,
			&e.ProjectID,
			&e.WeekStartDate,
			&e.WeekNumber,
			&e.PlannedHours,
			&e.ActualHours,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weekly entry: %w", err)
		}

		e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		// Determine if this is a past week (can edit actual hours)
		weekStart, _ := time.Parse("2006-01-02", e.WeekStartDate)
		isPastWeek := weekStart.Before(currentMonday)

		entryWithStatus := models.WeeklyEntryWithStatus{
			WeeklyEntry: e,
			IsPastWeek:  isPastWeek,
		}
		entryWithStatus.CalculateStatus()

		entries = append(entries, entryWithStatus)
	}

	if entries == nil {
		entries = []models.WeeklyEntryWithStatus{}
	}

	return entries, nil
}

// UpdateActualHours updates the actual hours for a specific weekly entry
func (a *App) UpdateActualHours(input models.UpdateActualHoursInput) (*models.WeeklyEntryWithStatus, error) {
	if input.ActualHours < 0 {
		return nil, models.ErrInvalidHours
	}

	// Get the entry to verify it exists and check if it's a past week
	var weekStartDate string
	err := a.db.QueryRow(`SELECT week_start_date FROM weekly_entries WHERE id = ?`, input.EntryID).Scan(&weekStartDate)
	if err != nil {
		return nil, models.ErrEntryNotFound
	}

	// Check if this is a past or current week (allow editing)
	weekStart, _ := time.Parse("2006-01-02", weekStartDate)
	currentMonday := getCurrentWeekMonday()
	nextMonday := currentMonday.AddDate(0, 0, 7)

	if weekStart.After(nextMonday) || weekStart.Equal(nextMonday) {
		return nil, models.ErrCannotEditFutureWeek
	}

	// Update the actual hours
	result, err := a.db.Exec(database.UpdateActualHours, input.ActualHours, input.EntryID)
	if err != nil {
		return nil, fmt.Errorf("failed to update actual hours: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, models.ErrEntryNotFound
	}

	// Return the updated entry
	return a.getWeeklyEntry(input.EntryID)
}

// GetWeeklyEntriesByWeek returns all weekly entries for a specific week across all active projects
func (a *App) GetWeeklyEntriesByWeek(weekStartDate string) ([]models.WeeklyEntryWithProject, error) {
	rows, err := a.db.Query(database.SelectWeeklyEntriesByWeek, weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query weekly entries by week: %w", err)
	}
	defer rows.Close()

	currentMonday := getCurrentWeekMonday()
	var entries []models.WeeklyEntryWithProject

	for rows.Next() {
		var e models.WeeklyEntry
		var projectName string
		var createdAt, updatedAt string

		err := rows.Scan(
			&e.ID,
			&e.ProjectID,
			&e.WeekStartDate,
			&e.WeekNumber,
			&e.PlannedHours,
			&e.ActualHours,
			&createdAt,
			&updatedAt,
			&projectName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weekly entry: %w", err)
		}

		e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		// Determine if this is a past week (can edit actual hours)
		weekStart, _ := time.Parse("2006-01-02", e.WeekStartDate)
		isPastWeek := weekStart.Before(currentMonday)

		entryWithStatus := models.WeeklyEntryWithStatus{
			WeeklyEntry: e,
			IsPastWeek:  isPastWeek,
		}
		entryWithStatus.CalculateStatus()

		entryWithProject := models.WeeklyEntryWithProject{
			WeeklyEntryWithStatus: entryWithStatus,
			ProjectName:           projectName,
		}

		entries = append(entries, entryWithProject)
	}

	if entries == nil {
		entries = []models.WeeklyEntryWithProject{}
	}

	return entries, nil
}

// getWeeklyEntry retrieves a single weekly entry by ID
func (a *App) getWeeklyEntry(entryID int64) (*models.WeeklyEntryWithStatus, error) {
	var e models.WeeklyEntry
	var createdAt, updatedAt string

	err := a.db.QueryRow(`
		SELECT id, project_id, week_start_date, week_number, planned_hours, actual_hours, created_at, updated_at
		FROM weekly_entries WHERE id = ?
	`, entryID).Scan(
		&e.ID,
		&e.ProjectID,
		&e.WeekStartDate,
		&e.WeekNumber,
		&e.PlannedHours,
		&e.ActualHours,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, models.ErrEntryNotFound
	}

	e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	currentMonday := getCurrentWeekMonday()
	weekStart, _ := time.Parse("2006-01-02", e.WeekStartDate)
	isPastWeek := weekStart.Before(currentMonday)

	entryWithStatus := &models.WeeklyEntryWithStatus{
		WeeklyEntry: e,
		IsPastWeek:  isPastWeek,
	}
	entryWithStatus.CalculateStatus()

	return entryWithStatus, nil
}
