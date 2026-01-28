package system

import (
	"fmt"
	"time"

	"btrack/internal/models"
	"btrack/internal/store"
)

// ProjectTemplate is an alias for store type
type ProjectTemplate = store.Template

// CreateProjectFromTemplateInput represents input for creating a project from a template
type CreateProjectFromTemplateInput struct {
	TemplateID int64  `json:"templateId"`
	Name       string `json:"name"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}

// CreateTemplate saves a project as a template
func (s *Service) CreateTemplate(projectID int64, templateName string) (*ProjectTemplate, error) {
	if templateName == "" {
		return nil, fmt.Errorf("template name is required")
	}

	// Get the project
	project, err := s.projectService.GetByID(s.ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Create template from project
	template := &store.Template{
		Name:            templateName,
		TotalSoldHours:  project.TotalSoldHours,
		SpecialistHours: project.SpecialistHours,
	}

	if err := s.store.CreateTemplate(template); err != nil {
		return nil, err
	}

	return template, nil
}

// GetTemplates returns all project templates
func (s *Service) GetTemplates() ([]ProjectTemplate, error) {
	return s.store.GetTemplates()
}

// GetTemplate returns a single template by ID
func (s *Service) GetTemplate(id int64) (*ProjectTemplate, error) {
	return s.store.GetTemplate(id)
}

// CreateProjectFromTemplate creates a new project from a template
func (s *Service) CreateProjectFromTemplate(input CreateProjectFromTemplateInput) (*models.ProjectWithStats, error) {
	// Get the template
	template, err := s.GetTemplate(input.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// Validate dates
	if input.StartDate == "" || input.EndDate == "" {
		return nil, fmt.Errorf("start date and end date are required")
	}

	if _, err := time.Parse("2006-01-02", input.StartDate); err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	if _, err := time.Parse("2006-01-02", input.EndDate); err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	// Create project using the project service
	createInput := models.CreateProjectInput{
		Name:            input.Name,
		TotalSoldHours:  template.TotalSoldHours,
		SpecialistHours: template.SpecialistHours,
		StartDate:       input.StartDate,
		EndDate:         input.EndDate,
	}

	return s.projectService.Create(s.ctx, createInput)
}

// DeleteTemplate deletes a template
func (s *Service) DeleteTemplate(id int64) error {
	return s.store.DeleteTemplate(id)
}
