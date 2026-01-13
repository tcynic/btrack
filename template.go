package main

import (
	"fmt"
	"time"

	"btrack/internal/database"
	"btrack/internal/models"
)

// ProjectTemplate represents a saved project template
type ProjectTemplate struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	TotalSoldHours  int       `json:"totalSoldHours"`
	SpecialistHours int       `json:"specialistHours"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// CreateTemplate saves a project as a template
func (a *App) CreateTemplate(projectID int64, templateName string) (*ProjectTemplate, error) {
	if templateName == "" {
		return nil, fmt.Errorf("template name is required")
	}

	// Get the project
	project, err := a.GetProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Create template from project
	result, err := a.db.Exec(database.InsertTemplate,
		templateName,
		project.TotalSoldHours,
		project.SpecialistHours,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert template: %w", err)
	}

	templateID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get template ID: %w", err)
	}

	return a.GetTemplate(templateID)
}

// GetTemplates returns all project templates
func (a *App) GetTemplates() ([]ProjectTemplate, error) {
	rows, err := a.db.Query(database.SelectAllTemplates)
	if err != nil {
		return nil, fmt.Errorf("failed to query templates: %w", err)
	}
	defer rows.Close()

	var templates []ProjectTemplate
	for rows.Next() {
		var t ProjectTemplate
		var createdAt, updatedAt string

		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.TotalSoldHours,
			&t.SpecialistHours,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		templates = append(templates, t)
	}

	if templates == nil {
		templates = []ProjectTemplate{}
	}

	return templates, nil
}

// GetTemplate returns a single template by ID
func (a *App) GetTemplate(id int64) (*ProjectTemplate, error) {
	var t ProjectTemplate
	var createdAt, updatedAt string

	err := a.db.QueryRow(database.SelectTemplateByID, id).Scan(
		&t.ID,
		&t.Name,
		&t.TotalSoldHours,
		&t.SpecialistHours,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}

	t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &t, nil
}

// CreateProjectFromTemplateInput represents input for creating a project from a template
type CreateProjectFromTemplateInput struct {
	TemplateID int64  `json:"templateId"`
	Name       string `json:"name"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}

// CreateProjectFromTemplate creates a new project from a template
func (a *App) CreateProjectFromTemplate(input CreateProjectFromTemplateInput) (*models.ProjectWithStats, error) {
	// Get the template
	template, err := a.GetTemplate(input.TemplateID)
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

	// Create project using the standard CreateProject method
	createInput := models.CreateProjectInput{
		Name:            input.Name,
		TotalSoldHours:  template.TotalSoldHours,
		SpecialistHours: template.SpecialistHours,
		StartDate:       input.StartDate,
		EndDate:         input.EndDate,
	}

	return a.CreateProject(createInput)
}

// DeleteTemplate deletes a template
func (a *App) DeleteTemplate(id int64) error {
	result, err := a.db.Exec(database.DeleteTemplate, id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}
