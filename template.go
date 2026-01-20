package main

import (
	"btrack/internal/models"
	"btrack/internal/services/system"
)

// ProjectTemplate represents a saved project template
type ProjectTemplate = system.ProjectTemplate

// CreateProjectFromTemplateInput represents input for creating a project from a template
type CreateProjectFromTemplateInput = system.CreateProjectFromTemplateInput

// CreateTemplate saves a project as a template
func (a *App) CreateTemplate(projectID int64, templateName string) (*ProjectTemplate, error) {
	return a.systemService.CreateTemplate(projectID, templateName)
}

// GetTemplates returns all project templates
func (a *App) GetTemplates() ([]ProjectTemplate, error) {
	return a.systemService.GetTemplates()
}

// GetTemplate returns a single template by ID
func (a *App) GetTemplate(id int64) (*ProjectTemplate, error) {
	return a.systemService.GetTemplate(id)
}

// CreateProjectFromTemplate creates a new project from a template
func (a *App) CreateProjectFromTemplate(input CreateProjectFromTemplateInput) (*models.ProjectWithStats, error) {
	return a.systemService.CreateProjectFromTemplate(input)
}

// DeleteTemplate deletes a template
func (a *App) DeleteTemplate(id int64) error {
	return a.systemService.DeleteTemplate(id)
}
