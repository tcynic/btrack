package main

import (
	"btrack/internal/models"
)

// CreateTask creates a new task
func (a *App) CreateTask(input models.CreateTaskInput) (*models.Task, error) {
	return a.notesService.CreateTask(input)
}

// GetTask returns a single task by ID
func (a *App) GetTask(id int64) (*models.Task, error) {
	return a.notesService.GetTask(id)
}

// GetTasksByProject returns all tasks for a project
func (a *App) GetTasksByProject(projectID int64) ([]models.Task, error) {
	return a.notesService.GetTasksByProject(projectID)
}

// GetTasksBySource returns tasks linked to a specific source (meeting or note)
func (a *App) GetTasksBySource(sourceType string, sourceID int64) ([]models.Task, error) {
	return a.notesService.GetTasksBySource(sourceType, sourceID)
}

// GetAllTasks returns all tasks across projects with optional filters
func (a *App) GetAllTasks(statusFilter string, projectIDFilter int64) ([]models.TaskWithContext, error) {
	return a.notesService.GetAllTasks(statusFilter, projectIDFilter)
}

// UpdateTask updates an existing task
func (a *App) UpdateTask(input models.UpdateTaskInput) (*models.Task, error) {
	return a.notesService.UpdateTask(input)
}

// UpdateTaskStatus updates only the status of a task
func (a *App) UpdateTaskStatus(id int64, status string) (*models.Task, error) {
	return a.notesService.UpdateTaskStatus(id, status)
}

// DeleteTask removes a task
func (a *App) DeleteTask(id int64) error {
	return a.notesService.DeleteTask(id)
}

// SearchTasks searches for tasks by title or description
func (a *App) SearchTasks(query string) ([]models.TaskWithContext, error) {
	return a.notesService.SearchTasks(query)
}
