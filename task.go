package main

import (
	"fmt"
	"strings"
	"time"

	"btrack/internal/database"
	"btrack/internal/models"
)

// CreateTask creates a new task
func (a *App) CreateTask(input models.CreateTaskInput) (*models.Task, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	result, err := a.db.Exec(database.InsertTask,
		input.ProjectID,
		input.SourceType,
		input.SourceID,
		input.Title,
		input.Description,
		input.Priority,
		input.DueDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get task ID: %w", err)
	}

	return a.GetTask(taskID)
}

// GetTask returns a single task by ID
func (a *App) GetTask(id int64) (*models.Task, error) {
	var t models.Task
	var sourceID *int64
	var description, dueDate *string
	var createdAt, updatedAt string

	err := a.db.QueryRow(database.SelectTaskByID, id).Scan(
		&t.ID,
		&t.ProjectID,
		&t.SourceType,
		&sourceID,
		&t.Title,
		&description,
		&t.Status,
		&t.Priority,
		&dueDate,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, models.ErrTaskNotFound
	}

	t.SourceID = sourceID
	if description != nil {
		t.Description = *description
	}
	if dueDate != nil {
		t.DueDate = *dueDate
	}
	t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

	return &t, nil
}

// GetTasksByProject returns all tasks for a project
func (a *App) GetTasksByProject(projectID int64) ([]models.Task, error) {
	rows, err := a.db.Query(database.SelectTasksByProject, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var sourceID *int64
		var description, dueDate *string
		var createdAt, updatedAt string

		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.SourceType,
			&sourceID,
			&t.Title,
			&description,
			&t.Status,
			&t.Priority,
			&dueDate,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		t.SourceID = sourceID
		if description != nil {
			t.Description = *description
		}
		if dueDate != nil {
			t.DueDate = *dueDate
		}
		t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

// GetTasksBySource returns tasks linked to a specific source (meeting or note)
func (a *App) GetTasksBySource(sourceType string, sourceID int64) ([]models.Task, error) {
	rows, err := a.db.Query(database.SelectTasksBySource, sourceType, sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var sID *int64
		var description, dueDate *string
		var createdAt, updatedAt string

		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.SourceType,
			&sID,
			&t.Title,
			&description,
			&t.Status,
			&t.Priority,
			&dueDate,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		t.SourceID = sID
		if description != nil {
			t.Description = *description
		}
		if dueDate != nil {
			t.DueDate = *dueDate
		}
		t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

// GetAllTasks returns all tasks across projects with optional filters
func (a *App) GetAllTasks(statusFilter string, projectIDFilter int64) ([]models.TaskWithContext, error) {
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

	rows, err := a.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.TaskWithContext
	for rows.Next() {
		var t models.TaskWithContext
		var sourceID *int64
		var description, dueDate, sourceTitle *string
		var createdAt, updatedAt string

		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.SourceType,
			&sourceID,
			&t.Title,
			&description,
			&t.Status,
			&t.Priority,
			&dueDate,
			&createdAt,
			&updatedAt,
			&t.ProjectName,
			&sourceTitle,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		t.SourceID = sourceID
		if description != nil {
			t.Description = *description
		}
		if dueDate != nil {
			t.DueDate = *dueDate
		}
		if sourceTitle != nil {
			t.SourceTitle = *sourceTitle
		}
		t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []models.TaskWithContext{}
	}

	return tasks, nil
}

// UpdateTask updates an existing task
func (a *App) UpdateTask(input models.UpdateTaskInput) (*models.Task, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	result, err := a.db.Exec(database.UpdateTask,
		input.Title,
		input.Description,
		input.Status,
		input.Priority,
		input.DueDate,
		input.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, models.ErrTaskNotFound
	}

	return a.GetTask(input.ID)
}

// UpdateTaskStatus updates only the status of a task
func (a *App) UpdateTaskStatus(id int64, status string) (*models.Task, error) {
	if !models.IsValidTaskStatus(status) {
		return nil, models.ErrInvalidStatus
	}

	result, err := a.db.Exec(database.UpdateTaskStatus, status, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, models.ErrTaskNotFound
	}

	return a.GetTask(id)
}

// DeleteTask removes a task
func (a *App) DeleteTask(id int64) error {
	result, err := a.db.Exec(database.DeleteTask, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrTaskNotFound
	}

	return nil
}

// SearchTasks searches for tasks by title or description
func (a *App) SearchTasks(query string) ([]models.TaskWithContext, error) {
	rows, err := a.db.Query(database.SearchTasks, query, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.TaskWithContext
	for rows.Next() {
		var t models.TaskWithContext
		var sourceID *int64
		var description, dueDate, sourceTitle *string
		var createdAt, updatedAt string

		err := rows.Scan(
			&t.ID,
			&t.ProjectID,
			&t.SourceType,
			&sourceID,
			&t.Title,
			&description,
			&t.Status,
			&t.Priority,
			&dueDate,
			&createdAt,
			&updatedAt,
			&t.ProjectName,
			&sourceTitle,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		t.SourceID = sourceID
		if description != nil {
			t.Description = *description
		}
		if dueDate != nil {
			t.DueDate = *dueDate
		}
		if sourceTitle != nil {
			t.SourceTitle = *sourceTitle
		}
		t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)

		tasks = append(tasks, t)
	}

	if tasks == nil {
		tasks = []models.TaskWithContext{}
	}

	return tasks, nil
}
