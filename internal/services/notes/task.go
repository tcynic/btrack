package notes

import (
	"fmt"
	"strings"

	"btrack/internal/database"
	"btrack/internal/models"
)

// CreateTask creates a new task
func (s *Service) CreateTask(input models.CreateTaskInput) (*models.Task, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	taskID, err := s.taskRepo.Create(
		input.ProjectID,
		input.SourceType,
		input.SourceID,
		input.Title,
		input.Description,
		input.Priority,
		input.DueDate,
	)
	if err != nil {
		return nil, err
	}

	return s.GetTask(taskID)
}

// GetTask returns a single task by ID
func (s *Service) GetTask(id int64) (*models.Task, error) {
	return s.taskRepo.GetByID(id)
}

// GetTasksByProject returns all tasks for a project
func (s *Service) GetTasksByProject(projectID int64) ([]models.Task, error) {
	tasks, err := s.taskRepo.GetByProject(projectID)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(tasks), nil
}

// GetTasksBySource returns tasks linked to a specific source (meeting or note)
func (s *Service) GetTasksBySource(sourceType string, sourceID int64) ([]models.Task, error) {
	tasks, err := s.taskRepo.GetBySource(sourceType, sourceID)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(tasks), nil
}

// GetAllTasks returns all tasks across projects with optional filters
func (s *Service) GetAllTasks(statusFilter string, projectIDFilter int64) ([]models.TaskWithContext, error) {
	query := database.SelectAllTasksFiltered
	args := []interface{}{}

	// Build WHERE clauses for filters
	var conditions []string
	if statusFilter != "" {
		conditions = append(conditions, "t.status = ?")
		args = append(args, statusFilter)
	}
	if projectIDFilter > 0 {
		conditions = append(conditions, "t.project_id = ?")
		args = append(args, projectIDFilter)
	}

	// Add filter conditions to query
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Add ordering
	query += ` ORDER BY CASE t.status
		WHEN 'in_progress' THEN 1
		WHEN 'pending' THEN 2
		WHEN 'completed' THEN 3
		WHEN 'cancelled' THEN 4
	END, t.priority DESC, t.due_date ASC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.TaskWithContext
	for rows.Next() {
		t, err := models.ScanTaskWithContext(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, *t)
	}

	return database.EnsureSlice(tasks), nil
}

// UpdateTask updates an existing task
func (s *Service) UpdateTask(input models.UpdateTaskInput) (*models.Task, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	err := s.taskRepo.Update(
		input.ID,
		input.Title,
		input.Description,
		input.Status,
		input.Priority,
		input.DueDate,
	)
	if err != nil {
		return nil, err
	}

	return s.GetTask(input.ID)
}

// UpdateTaskStatus updates only the status of a task
func (s *Service) UpdateTaskStatus(id int64, status string) (*models.Task, error) {
	if !models.IsValidTaskStatus(status) {
		return nil, models.ValidationError("status", "invalid status value")
	}

	err := s.taskRepo.UpdateStatus(id, status)
	if err != nil {
		return nil, err
	}

	return s.GetTask(id)
}

// DeleteTask removes a task
func (s *Service) DeleteTask(id int64) error {
	return s.taskRepo.Delete(id)
}

// SearchTasks searches for tasks by title or description
func (s *Service) SearchTasks(query string) ([]models.TaskWithContext, error) {
	searchPattern := "%" + query + "%"
	tasks, err := s.taskRepo.Search(searchPattern)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(tasks), nil
}
