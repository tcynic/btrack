package repository

import (
	"btrack/internal/database"
	"btrack/internal/models"
)

// ProjectRepository handles project-related database operations
type ProjectRepository struct {
	*Repository
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(base *Repository) *ProjectRepository {
	return &ProjectRepository{Repository: base}
}

// GetByID retrieves a single project by ID
func (r *ProjectRepository) GetByID(id int64) (*models.Project, error) {
	return r.QueryOne(
		database.SelectProjectByID,
		[]any{id},
		models.ScanProject,
	)
}

// GetAll retrieves all projects, optionally filtered by active status
func (r *ProjectRepository) GetAll(activeOnly bool) ([]models.Project, error) {
	query := database.SelectAllProjects
	if activeOnly {
		query = database.SelectActiveProjects
	}
	return r.QuerySlice(query, nil, models.ScanProject)
}

// Search searches for projects by name
func (r *ProjectRepository) Search(searchQuery string) ([]models.Project, error) {
	return r.QuerySlice(
		database.SearchProjects,
		[]any{searchQuery},
		models.ScanProject,
	)
}

// Create inserts a new project and returns the ID
func (r *ProjectRepository) Create(name string, totalSoldHours, specialistHours int, startDate, endDate string) (int64, error) {
	result, err := r.Exec(
		database.InsertProject,
		name, totalSoldHours, specialistHours, startDate, endDate,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update updates an existing project
func (r *ProjectRepository) Update(id int64, name string, totalSoldHours, specialistHours int, startDate, endDate string, isActive bool) error {
	activeInt := 0
	if isActive {
		activeInt = 1
	}
	
	result, err := r.Exec(
		database.UpdateProject,
		name, totalSoldHours, specialistHours, startDate, endDate, activeInt, id,
	)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("project")
	}
	
	return nil
}

// SoftDelete marks a project as deleted
func (r *ProjectRepository) SoftDelete(id int64) error {
	result, err := r.Exec(database.SoftDeleteProject, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("project")
	}
	
	return nil
}

// Restore restores a soft-deleted project
func (r *ProjectRepository) Restore(id int64) error {
	result, err := r.Exec(database.RestoreProject, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("project")
	}
	
	return nil
}

// PermanentlyDelete permanently removes a project
func (r *ProjectRepository) PermanentlyDelete(id int64) error {
	result, err := r.Exec(database.PermanentlyDeleteProject, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("project")
	}
	
	return nil
}

// GetStats retrieves aggregated statistics for a project
func (r *ProjectRepository) GetStats(projectID int64) (totalPlanned, totalActual, totalWeeks int, err error) {
	err = r.DB().QueryRow(database.SelectProjectStats, projectID).Scan(
		&totalPlanned, &totalActual, &totalWeeks,
	)
	if err != nil {
		return 0, 0, 0, models.DatabaseError(err)
	}
	return totalPlanned, totalActual, totalWeeks, nil
}
