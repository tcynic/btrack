package main

import (
	"fmt"

	"btrack/internal/database"
	"btrack/internal/models"
)

// CreateGoal creates a new goal for a project
func (a *App) CreateGoal(input models.CreateGoalInput) (*models.Goal, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	result, err := a.db.Exec(database.InsertGoal,
		input.ProjectID,
		input.Title,
		input.Description,
		input.TargetDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert goal: %w", err)
	}

	goalID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get goal ID: %w", err)
	}

	return a.GetGoal(goalID)
}

// GetGoals returns all goals for a project
func (a *App) GetGoals(projectID int64) ([]models.Goal, error) {
	rows, err := a.db.Query(database.SelectGoalsByProject, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query goals: %w", err)
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		g, err := models.ScanGoal(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, *g)
	}

	return database.EnsureSlice(goals), nil
}

// GetGoal returns a single goal by ID
func (a *App) GetGoal(id int64) (*models.Goal, error) {
	row := a.db.QueryRow(database.SelectGoalByID, id)
	g, err := models.ScanGoal(row.Scan)
	if err != nil {
		return nil, models.ErrGoalNotFound
	}
	return g, nil
}

// UpdateGoal updates an existing goal
func (a *App) UpdateGoal(input models.UpdateGoalInput) (*models.Goal, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Default status if not provided
	if input.Status == "" {
		input.Status = models.GoalStatusPending
	}

	result, err := a.db.Exec(database.UpdateGoal,
		input.Title,
		input.Description,
		input.Status,
		input.TargetDate,
		input.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, models.ErrGoalNotFound
	}

	return a.GetGoal(input.ID)
}

// UpdateGoalStatus updates only the status of a goal
func (a *App) UpdateGoalStatus(id int64, status string) (*models.Goal, error) {
	if !models.IsValidStatus(status) {
		return nil, models.ErrInvalidStatus
	}

	result, err := a.db.Exec(database.UpdateGoalStatus, status, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update goal status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, models.ErrGoalNotFound
	}

	return a.GetGoal(id)
}

// DeleteGoal removes a goal
func (a *App) DeleteGoal(id int64) error {
	result, err := a.db.Exec(database.DeleteGoal, id)
	if err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrGoalNotFound
	}

	return nil
}

// GoalStats represents statistics for a project's goals
type GoalStats struct {
	Total          int     `json:"total"`
	Pending        int     `json:"pending"`
	InProgress     int     `json:"inProgress"`
	Completed      int     `json:"completed"`
	Cancelled      int     `json:"cancelled"`
	CompletionRate float64 `json:"completionRate"` // percentage
}

// GetGoalStats returns statistics for a project's goals
func (a *App) GetGoalStats(projectID int64) (*GoalStats, error) {
	rows, err := a.db.Query(database.SelectGoalStatsByProject, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query goal stats: %w", err)
	}
	defer rows.Close()

	stats := &GoalStats{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan goal stat: %w", err)
		}

		stats.Total += count
		switch status {
		case models.GoalStatusPending:
			stats.Pending = count
		case models.GoalStatusInProgress:
			stats.InProgress = count
		case models.GoalStatusCompleted:
			stats.Completed = count
		case models.GoalStatusCancelled:
			stats.Cancelled = count
		}
	}

	// Calculate completion rate
	if stats.Total > 0 {
		stats.CompletionRate = float64(stats.Completed) / float64(stats.Total) * 100
	}

	return stats, nil
}

// SearchGoals searches for goals by title or description
func (a *App) SearchGoals(query string) ([]models.GoalWithProject, error) {
	rows, err := a.db.Query(database.SearchGoals, query, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search goals: %w", err)
	}
	defer rows.Close()

	var goals []models.GoalWithProject
	for rows.Next() {
		g, err := models.ScanGoalWithProject(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}
		goals = append(goals, *g)
	}

	return database.EnsureSlice(goals), nil
}
