package main

import (
	"fmt"
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
		var p models.Project
		var isActive int
		var createdAt, updatedAt string

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.TotalSoldHours,
			&p.SpecialistHours,
			&p.StartDate,
			&p.EndDate,
			&isActive,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		p.IsActive = isActive == 1
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		// Get stats for this project
		stats, err := a.getProjectStats(p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get project stats: %w", err)
		}

		projects = append(projects, models.ProjectWithStats{
			Project:           p,
			MyHours:           p.TotalSoldHours - p.SpecialistHours,
			TotalWeeks:        stats.TotalWeeks,
			TotalPlannedHours: stats.TotalPlanned,
			TotalActualHours:  stats.TotalActual,
		})
	}

	if projects == nil {
		projects = []models.ProjectWithStats{}
	}

	return projects, nil
}

// GetProject returns a single project by ID with stats
func (a *App) GetProject(id int64) (*models.ProjectWithStats, error) {
	var p models.Project
	var isActive int
	var createdAt, updatedAt string

	err := a.db.QueryRow(database.SelectProjectByID, id).Scan(
		&p.ID,
		&p.Name,
		&p.TotalSoldHours,
		&p.SpecialistHours,
		&p.StartDate,
		&p.EndDate,
		&isActive,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, models.ErrProjectNotFound
	}

	p.IsActive = isActive == 1
	p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	stats, err := a.getProjectStats(p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project stats: %w", err)
	}

	return &models.ProjectWithStats{
		Project:           p,
		MyHours:           p.TotalSoldHours - p.SpecialistHours,
		TotalWeeks:        stats.TotalWeeks,
		TotalPlannedHours: stats.TotalPlanned,
		TotalActualHours:  stats.TotalActual,
	}, nil
}

// UpdateProject updates project details and recalculates weekly entries if needed
func (a *App) UpdateProject(input models.UpdateProjectInput) (*models.ProjectWithStats, error) {
	// Verify project exists
	existing, err := a.GetProject(input.ID)
	if err != nil {
		return nil, err
	}

	// Check if we need to recalculate entries
	needsRecalc := existing.TotalSoldHours != input.TotalSoldHours ||
		existing.SpecialistHours != input.SpecialistHours ||
		existing.StartDate != input.StartDate ||
		existing.EndDate != input.EndDate

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
		// Delete existing entries
		_, err = tx.Exec(database.DeleteWeeklyEntriesByProject, input.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete weekly entries: %w", err)
		}

		// Calculate new distribution
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
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return a.GetProject(input.ID)
}

// DeleteProject removes a project and its weekly entries
func (a *App) DeleteProject(id int64) error {
	result, err := a.db.Exec(database.DeleteProject, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrProjectNotFound
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
