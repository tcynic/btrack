package main

import (
	"fmt"
	"log"
	"time"

	"btrack/internal/database"
	"btrack/internal/models"
)

// CreateProject creates a new project and generates weekly entries based on frontloading
func (a *App) CreateProject(input models.CreateProjectInput) (*models.ProjectWithStats, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Calculate the weekly distribution
	entries, err := a.CalculateDistribution(input)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate distribution: %w", err)
	}

	// Start transaction
	tx, err := a.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert project
	result, err := tx.Exec(database.InsertProject,
		input.Name,
		input.TotalSoldHours,
		input.SpecialistHours,
		input.StartDate,
		input.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert project: %w", err)
	}

	projectID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get project ID: %w", err)
	}

	// Insert weekly entries
	for _, entry := range entries {
		_, err := tx.Exec(database.InsertWeeklyEntry,
			projectID,
			entry.WeekStartDate,
			entry.WeekNumber,
			entry.PlannedHours,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert weekly entry: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return the created project with stats
	return a.GetProject(projectID)
}

// GetAllProjects returns all projects, optionally filtered by active status
func (a *App) GetAllProjects(activeOnly bool) ([]models.ProjectWithStats, error) {
	query := database.SelectAllProjects
	if activeOnly {
		query = database.SelectActiveProjects
	}

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []models.ProjectWithStats
	for rows.Next() {
		p, err := models.ScanProject(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		// Get stats for this project
		stats, err := a.getProjectStats(p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get project stats: %w", err)
		}

		projectWithStats := models.ProjectWithStats{
			Project:           *p,
			MyHours:           p.TotalSoldHours - p.SpecialistHours,
			TotalWeeks:        stats.TotalWeeks,
			TotalPlannedHours: stats.TotalPlanned,
			TotalActualHours:  stats.TotalActual,
		}
		projectWithStats.Health = calculateProjectHealth(projectWithStats)
		projects = append(projects, projectWithStats)
	}

	if projects == nil {
		projects = []models.ProjectWithStats{}
	}

	return projects, nil
}

// GetProject returns a single project by ID with stats
func (a *App) GetProject(id int64) (*models.ProjectWithStats, error) {
	row := a.db.QueryRow(database.SelectProjectByID, id)
	p, err := models.ScanProject(row.Scan)
	if err != nil {
		return nil, models.NotFound("project")
	}

	stats, err := a.getProjectStats(p.ID)
	if err != nil {
		return nil, models.Internal(err, "failed to get project stats")
	}

	projectWithStats := &models.ProjectWithStats{
		Project:           *p,
		MyHours:           p.TotalSoldHours - p.SpecialistHours,
		TotalWeeks:        stats.TotalWeeks,
		TotalPlannedHours: stats.TotalPlanned,
		TotalActualHours:  stats.TotalActual,
	}
	projectWithStats.Health = calculateProjectHealth(*projectWithStats)
	return projectWithStats, nil
}

// UpdateProject updates project details and recalculates weekly entries if needed
func (a *App) UpdateProject(input models.UpdateProjectInput) (*models.ProjectWithStats, error) {
	// Verify project exists
	existing, err := a.GetProject(input.ID)
	if err != nil {
		return nil, err
	}

	// Check if we need to recalculate entries (skip for persistent projects)
	needsRecalc := !existing.IsPersistent && (existing.TotalSoldHours != input.TotalSoldHours ||
		existing.SpecialistHours != input.SpecialistHours ||
		existing.StartDate != input.StartDate ||
		existing.EndDate != input.EndDate)

	tx, err := a.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update project
	isActive := 0
	if input.IsActive {
		isActive = 1
	}

	_, err = tx.Exec(database.UpdateProject,
		input.Name,
		input.TotalSoldHours,
		input.SpecialistHours,
		input.StartDate,
		input.EndDate,
		isActive,
		input.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Recalculate entries if needed
	if needsRecalc {
		// Get next Monday (start of future weeks)
		nextMonday := getCurrentWeekMonday().AddDate(0, 0, 7)
		nextMondayStr := nextMonday.Format("2006-01-02")

		// Parse the new start and end dates
		newStartDate, err := time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}

		newEndDate, err := time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}

		// If the new start date is after next Monday, we need to delete everything and recalculate
		if newStartDate.After(nextMonday) || newStartDate.Equal(nextMonday) {
			// Delete all existing entries
			_, err = tx.Exec(database.DeleteWeeklyEntriesByProject, input.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to delete weekly entries: %w", err)
			}

			// Calculate new distribution from scratch
			entries, err := a.CalculateDistribution(models.CreateProjectInput{
				Name:            input.Name,
				TotalSoldHours:  input.TotalSoldHours,
				SpecialistHours: input.SpecialistHours,
				StartDate:       input.StartDate,
				EndDate:         input.EndDate,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to calculate distribution: %w", err)
			}

			// Insert new entries
			for _, entry := range entries {
				_, err := tx.Exec(database.InsertWeeklyEntry,
					input.ID,
					entry.WeekStartDate,
					entry.WeekNumber,
					entry.PlannedHours,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to insert weekly entry: %w", err)
				}
			}
		} else {
			// Preserve past/current weeks, only recalculate future weeks

			// Get sum of planned hours for past/current weeks
			var pastWeekHours int
			err = tx.QueryRow(database.SelectPastWeeklyHours, input.ID, nextMondayStr).Scan(&pastWeekHours)
			if err != nil {
				return nil, fmt.Errorf("failed to get past week hours: %w", err)
			}

			// Delete only future week entries
			_, err = tx.Exec(database.DeleteFutureWeeklyEntries, input.ID, nextMondayStr)
			if err != nil {
				return nil, fmt.Errorf("failed to delete future weekly entries: %w", err)
			}

			// Calculate remaining hours and weeks
			myHours := input.TotalSoldHours - input.SpecialistHours
			remainingHours := myHours - pastWeekHours

			// Only create future entries if there are remaining hours and the end date is in the future
			if remainingHours > 0 && newEndDate.After(nextMonday) {
				// Calculate weeks between next Monday and end date
				futureWeeks := calculateWeeksBetween(nextMonday, newEndDate)
				if futureWeeks < 1 {
					futureWeeks = 1
				}

				// Distribute remaining hours across future weeks (frontloaded)
				baseHours := remainingHours / futureWeeks
				remainder := remainingHours % futureWeeks

				currentMonday := nextMonday
				for i := 0; i < futureWeeks; i++ {
					hours := baseHours
					if i < remainder {
						hours++ // Frontload extra hours
					}

					_, err := tx.Exec(database.InsertWeeklyEntry,
						input.ID,
						currentMonday.Format("2006-01-02"),
						0, // Week number will be recalculated; set to 0 for now
						hours,
					)
					if err != nil {
						return nil, fmt.Errorf("failed to insert weekly entry: %w", err)
					}

					currentMonday = currentMonday.AddDate(0, 0, 7)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return a.GetProject(input.ID)
}

// DeleteProject soft-deletes a project (sets deleted_at timestamp)
func (a *App) DeleteProject(id int64) error {
	// Check if project is persistent
	project, err := a.GetProject(id)
	if err != nil {
		return err
	}

	if project.IsPersistent {
		return models.Forbidden("cannot delete persistent project")
	}

	// Auto-backup before deletion
	if err := a.autoBackup(); err != nil {
		// Log warning but don't fail the operation
		log.Printf("Warning: auto-backup failed before delete: %v", err)
	}

	result, err := a.db.Exec(database.SoftDeleteProject, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.NotFound("project")
	}

	return nil
}

// RestoreProject restores a soft-deleted project
func (a *App) RestoreProject(id int64) (*models.ProjectWithStats, error) {
	result, err := a.db.Exec(database.RestoreProject, id)
	if err != nil {
		return nil, fmt.Errorf("failed to restore project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, models.NotFound("project")
	}

	// Return the restored project
	// Note: We need to query without the deleted_at filter
	row := a.db.QueryRow(`
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, is_persistent, created_at, updated_at
		FROM projects WHERE id = ?
	`, id)
	p, err := models.ScanProject(row.Scan)
	if err != nil {
		return nil, models.NotFound("project")
	}

	stats, err := a.getProjectStats(p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project stats: %w", err)
	}

	return &models.ProjectWithStats{
		Project:           *p,
		MyHours:           p.TotalSoldHours - p.SpecialistHours,
		TotalWeeks:        stats.TotalWeeks,
		TotalPlannedHours: stats.TotalPlanned,
		TotalActualHours:  stats.TotalActual,
	}, nil
}

// PermanentlyDeleteProject permanently removes a project and its weekly entries
func (a *App) PermanentlyDeleteProject(id int64) error {
	// Check if project is persistent
	// Query without deleted_at filter to check persistence
	var isPersistent int
	err := a.db.QueryRow(`
		SELECT is_persistent FROM projects WHERE id = ?
	`, id).Scan(&isPersistent)
	if err != nil {
		return models.NotFound("project")
	}

	if isPersistent == 1 {
		return models.Forbidden("cannot permanently delete persistent project")
	}

	// Auto-backup before permanent deletion
	if err := a.autoBackup(); err != nil {
		// Log warning but don't fail the operation
		log.Printf("Warning: auto-backup failed before permanent delete: %v", err)
	}

	result, err := a.db.Exec(database.PermanentlyDeleteProject, id)
	if err != nil {
		return fmt.Errorf("failed to permanently delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.NotFound("project")
	}

	return nil
}

// ToggleProjectActive toggles the is_active status of a project
func (a *App) ToggleProjectActive(id int64) (*models.ProjectWithStats, error) {
	project, err := a.GetProject(id)
	if err != nil {
		return nil, err
	}

	return a.UpdateProject(models.UpdateProjectInput{
		ID:              project.ID,
		Name:            project.Name,
		TotalSoldHours:  project.TotalSoldHours,
		SpecialistHours: project.SpecialistHours,
		StartDate:       project.StartDate,
		EndDate:         project.EndDate,
		IsActive:        !project.IsActive,
	})
}

// projectStats holds aggregated stats for a project
type projectStats struct {
	TotalPlanned int
	TotalActual  int
	TotalWeeks   int
}

// SearchProjects searches for projects by name
func (a *App) SearchProjects(query string) ([]models.ProjectWithStats, error) {
	rows, err := a.db.Query(database.SearchProjects, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search projects: %w", err)
	}
	defer rows.Close()

	var projects []models.ProjectWithStats
	for rows.Next() {
		p, err := models.ScanProject(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		stats, err := a.getProjectStats(p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get project stats: %w", err)
		}

		projectWithStats := models.ProjectWithStats{
			Project:           *p,
			MyHours:           p.TotalSoldHours - p.SpecialistHours,
			TotalWeeks:        stats.TotalWeeks,
			TotalPlannedHours: stats.TotalPlanned,
			TotalActualHours:  stats.TotalActual,
		}
		projectWithStats.Health = calculateProjectHealth(projectWithStats)
		projects = append(projects, projectWithStats)
	}

	if projects == nil {
		projects = []models.ProjectWithStats{}
	}

	return projects, nil
}

// calculateProjectHealth determines the health status of a project
func calculateProjectHealth(p models.ProjectWithStats) models.ProjectHealth {
	// Persistent projects are always on track
	if p.IsPersistent {
		return models.ProjectHealth{
			Status:   models.HealthOnTrack,
			Message:  "Ongoing work",
			Severity: "info",
		}
	}

	// Parse end date
	endDate, err := time.Parse("2006-01-02", p.EndDate)
	if err != nil {
		// If we can't parse the date, default to on_track
		return models.ProjectHealth{
			Status:   models.HealthOnTrack,
			Message:  "Project status unavailable",
			Severity: "info",
		}
	}

	today := time.Now()
	isPastEnd := today.After(endDate)
	hoursRemaining := p.TotalPlannedHours - p.TotalActualHours
	variance := float64(p.TotalActualHours) / float64(p.TotalPlannedHours)

	// Completed: past end date
	if isPastEnd {
		if hoursRemaining > 0 {
			return models.ProjectHealth{
				Status:   models.HealthCompleted,
				Message:  fmt.Sprintf("Project ended with %d hours remaining", hoursRemaining),
				Severity: "info",
			}
		}
		return models.ProjectHealth{
			Status:   models.HealthCompleted,
			Message:  "Project completed",
			Severity: "info",
		}
	}

	// Over budget: actual > planned
	if p.TotalActualHours > p.TotalPlannedHours {
		overBy := p.TotalActualHours - p.TotalPlannedHours
		return models.ProjectHealth{
			Status:   models.HealthOverBudget,
			Message:  fmt.Sprintf("Over budget by %d hours", overBy),
			Severity: "error",
		}
	}

	// At risk: actual > 80% of planned
	if variance > 0.8 && hoursRemaining > 0 {
		return models.ProjectHealth{
			Status:   models.HealthAtRisk,
			Message:  fmt.Sprintf("%.0f%% capacity used, %d hours remaining", variance*100, hoursRemaining),
			Severity: "warning",
		}
	}

	// On track: everything normal
	return models.ProjectHealth{
		Status:   models.HealthOnTrack,
		Message:  fmt.Sprintf("%d hours remaining", hoursRemaining),
		Severity: "info",
	}
}

// getProjectStats retrieves aggregated stats for a project
func (a *App) getProjectStats(projectID int64) (*projectStats, error) {
	var stats projectStats
	err := a.db.QueryRow(database.SelectProjectStats, projectID).Scan(
		&stats.TotalPlanned,
		&stats.TotalActual,
		&stats.TotalWeeks,
	)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
