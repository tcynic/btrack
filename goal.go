package main

import (
	"fmt"
	"time"

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
		var g models.Goal
		var description, targetDate *string
		var createdAt, updatedAt string

		err := rows.Scan(
			&g.ID,
			&g.ProjectID,
			&g.Title,
			&description,
			&g.Status,
			&targetDate,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}

		if description != nil {
			g.Description = *description
		}
		if targetDate != nil {
			g.TargetDate = *targetDate
		}
		g.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		g.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		goals = append(goals, g)
	}

	if goals == nil {
		goals = []models.Goal{}
	}

	return goals, nil
}

// GetGoal returns a single goal by ID
func (a *App) GetGoal(id int64) (*models.Goal, error) {
	var g models.Goal
	var description, targetDate *string
	var createdAt, updatedAt string

	err := a.db.QueryRow(database.SelectGoalByID, id).Scan(
		&g.ID,
		&g.ProjectID,
		&g.Title,
		&description,
		&g.Status,
		&targetDate,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, models.ErrGoalNotFound
	}

	if description != nil {
		g.Description = *description
	}
	if targetDate != nil {
		g.TargetDate = *targetDate
	}
	g.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	g.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &g, nil
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
