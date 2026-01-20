package main

import (
	"btrack/internal/models"
	"btrack/internal/services/notes"
)

// CreateGoal creates a new goal for a project
func (a *App) CreateGoal(input models.CreateGoalInput) (*models.Goal, error) {
	return a.notesService.CreateGoal(input)
}

// GetGoals returns all goals for a project
func (a *App) GetGoals(projectID int64) ([]models.Goal, error) {
	return a.notesService.GetGoals(projectID)
}

// GetGoal returns a single goal by ID
func (a *App) GetGoal(id int64) (*models.Goal, error) {
	return a.notesService.GetGoal(id)
}

// UpdateGoal updates an existing goal
func (a *App) UpdateGoal(input models.UpdateGoalInput) (*models.Goal, error) {
	return a.notesService.UpdateGoal(input)
}

// UpdateGoalStatus updates only the status of a goal
func (a *App) UpdateGoalStatus(id int64, status string) (*models.Goal, error) {
	return a.notesService.UpdateGoalStatus(id, status)
}

// DeleteGoal removes a goal
func (a *App) DeleteGoal(id int64) error {
	return a.notesService.DeleteGoal(id)
}

// GoalStats represents statistics for a project's goals
type GoalStats = notes.GoalStats

// GetGoalStats returns statistics for a project's goals
func (a *App) GetGoalStats(projectID int64) (*GoalStats, error) {
	return a.notesService.GetGoalStats(projectID)
}

// SearchGoals searches for goals by title or description
func (a *App) SearchGoals(query string) ([]models.GoalWithProject, error) {
	return a.notesService.SearchGoals(query)
}
