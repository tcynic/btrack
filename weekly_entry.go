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
		e, err := models.ScanWeeklyEntry(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weekly entry: %w", err)
		}

		// Determine if this is a past week (can edit actual hours)
		weekStart, _ := time.Parse("2006-01-02", e.WeekStartDate)
		isPastWeek := weekStart.Before(currentMonday)

		entryWithStatus := models.WeeklyEntryWithStatus{
			WeeklyEntry: *e,
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
	// First, ensure all persistent projects have entries for this week
	if err := a.ensurePersistentEntries(weekStartDate); err != nil {
		return nil, fmt.Errorf("failed to ensure persistent entries: %w", err)
	}

	rows, err := a.db.Query(database.SelectWeeklyEntriesByWeek, weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query weekly entries by week: %w", err)
	}
	defer rows.Close()

	currentMonday := getCurrentWeekMonday()
	var entries []models.WeeklyEntryWithProject

	for rows.Next() {
		// Determine if this is a past week (can edit actual hours)
		// We need to parse the week start date to determine this
		var weekStartTemp string
		var tempRow struct {
			id, projectID, weekNumber, plannedHours, actualHours int64
			createdAt, updatedAt, projectName                      string
		}
		err := rows.Scan(
			&tempRow.id,
			&tempRow.projectID,
			&weekStartTemp,
			&tempRow.weekNumber,
			&tempRow.plannedHours,
			&tempRow.actualHours,
			&tempRow.createdAt,
			&tempRow.updatedAt,
			&tempRow.projectName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weekly entry: %w", err)
		}

		weekStart, _ := time.Parse("2006-01-02", weekStartTemp)
		isPastWeek := weekStart.Before(currentMonday)

		// Re-scan using the proper scanner
		// Note: This is a workaround. Ideally we'd refactor the scanner to accept isPastWeek calculation
		e := models.WeeklyEntryWithProject{
			WeeklyEntryWithStatus: models.WeeklyEntryWithStatus{
				WeeklyEntry: models.WeeklyEntry{
					ID:            tempRow.id,
					ProjectID:     tempRow.projectID,
					WeekStartDate: weekStartTemp,
					WeekNumber:    int(tempRow.weekNumber),
					PlannedHours:  int(tempRow.plannedHours),
					ActualHours:   int(tempRow.actualHours),
					CreatedAt:     database.ParseTimestamp(tempRow.createdAt),
					UpdatedAt:     database.ParseTimestamp(tempRow.updatedAt),
				},
				IsPastWeek: isPastWeek,
			},
			ProjectName: tempRow.projectName,
		}
		e.CalculateStatus()

		entryWithProject := e

		entries = append(entries, entryWithProject)
	}

	if entries == nil {
		entries = []models.WeeklyEntryWithProject{}
	}

	return entries, nil
}

// ensurePersistentEntries creates weekly entries for all persistent projects if they don't exist
func (a *App) ensurePersistentEntries(weekStartDate string) error {
	// Get all active persistent projects
	rows, err := a.db.Query("SELECT id FROM projects WHERE is_persistent = 1 AND is_active = 1")
	if err != nil {
		return fmt.Errorf("failed to query persistent projects: %w", err)
	}
	defer rows.Close()

	var projectIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan project ID: %w", err)
		}
		projectIDs = append(projectIDs, id)
	}

	// For each persistent project, insert entry if it doesn't exist
	for _, projectID := range projectIDs {
		_, err := a.db.Exec(`
			INSERT OR IGNORE INTO weekly_entries (project_id, week_start_date, week_number, planned_hours, actual_hours)
			VALUES (?, ?, 0, 0, 0)
		`, projectID, weekStartDate)
		if err != nil {
			return fmt.Errorf("failed to insert persistent entry: %w", err)
		}
	}

	return nil
}

// getWeeklyEntry retrieves a single weekly entry by ID
func (a *App) getWeeklyEntry(entryID int64) (*models.WeeklyEntryWithStatus, error) {
	row := a.db.QueryRow(`
		SELECT id, project_id, week_start_date, week_number, planned_hours, actual_hours, created_at, updated_at
		FROM weekly_entries WHERE id = ?
	`, entryID)
	e, err := models.ScanWeeklyEntry(row.Scan)
	if err != nil {
		return nil, models.ErrEntryNotFound
	}

	currentMonday := getCurrentWeekMonday()
	weekStart, _ := time.Parse("2006-01-02", e.WeekStartDate)
	isPastWeek := weekStart.Before(currentMonday)

	entryWithStatus := &models.WeeklyEntryWithStatus{
		WeeklyEntry: *e,
		IsPastWeek:  isPastWeek,
	}
	entryWithStatus.CalculateStatus()

	return entryWithStatus, nil
}
