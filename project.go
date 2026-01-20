package main

import (
	"btrack/internal/models"
)

// CreateProject creates a new project and generates weekly entries based on frontloading
func (a *App) CreateProject(input models.CreateProjectInput) (*models.ProjectWithStats, error) {
	return a.projectService.Create(a.ctx, input)
}

// GetAllProjects returns all projects, optionally filtered by active status
func (a *App) GetAllProjects(activeOnly bool) ([]models.ProjectWithStats, error) {
	return a.projectService.GetAll(a.ctx, activeOnly)
}

// GetProject returns a single project by ID with stats
func (a *App) GetProject(id int64) (*models.ProjectWithStats, error) {
	return a.projectService.GetByID(a.ctx, id)
}

// UpdateProject updates project details and recalculates weekly entries if needed
func (a *App) UpdateProject(input models.UpdateProjectInput) (*models.ProjectWithStats, error) {
	return a.projectService.Update(a.ctx, input, a.autoBackup)
}

// DeleteProject soft-deletes a project (sets deleted_at timestamp)
func (a *App) DeleteProject(id int64) error {
	return a.projectService.Delete(a.ctx, id, a.autoBackup)
}

// RestoreProject restores a soft-deleted project
func (a *App) RestoreProject(id int64) (*models.ProjectWithStats, error) {
	return a.projectService.Restore(a.ctx, id)
}

// PermanentlyDeleteProject permanently removes a project and its weekly entries
func (a *App) PermanentlyDeleteProject(id int64) error {
	return a.projectService.PermanentlyDelete(a.ctx, id, a.autoBackup)
}

// ToggleProjectActive toggles the is_active status of a project
func (a *App) ToggleProjectActive(id int64) (*models.ProjectWithStats, error) {
	return a.projectService.ToggleActive(a.ctx, id, a.autoBackup)
}

// SearchProjects searches for projects by name
func (a *App) SearchProjects(query string) ([]models.ProjectWithStats, error) {
	// Get basic projects from service
	projects, err := a.projectService.Search(a.ctx, query)
	if err != nil {
		return nil, err
	}

	// Convert to ProjectWithStats
	var result []models.ProjectWithStats
	for _, p := range projects {
		// Get full project with stats
		ps, err := a.GetProject(p.ID)
		if err != nil {
			continue // Skip projects we can't fetch
		}
		result = append(result, *ps)
	}

	if result == nil {
		result = []models.ProjectWithStats{}
	}

	return result, nil
}

