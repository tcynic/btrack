package repository

import (
	"btrack/internal/database"
	"btrack/internal/models"
)

// TaskRepository handles task database operations
type TaskRepository struct {
	*Repository
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(base *Repository) *TaskRepository {
	return &TaskRepository{Repository: base}
}

// GetByID retrieves a single task by ID
func (r *TaskRepository) GetByID(id int64) (*models.Task, error) {
	var task *models.Task
	err := r.QueryOne(database.SelectTaskByID, []any{id}, func(scan func(dest ...any) error) error {
		t, err := models.ScanTask(scan)
		if err != nil {
			return err
		}
		task = t
		return nil
	})
	if err != nil {
		return nil, err
	}
	return task, nil
}

// GetByProject retrieves all tasks for a project
func (r *TaskRepository) GetByProject(projectID int64) ([]models.Task, error) {
	var tasks []models.Task
	err := r.QuerySlice(database.SelectTasksByProject, []any{projectID}, func(scan func(dest ...any) error) error {
		t, err := models.ScanTask(scan)
		if err != nil {
			return err
		}
		tasks = append(tasks, *t)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}

// GetBySource retrieves tasks linked to a specific source
func (r *TaskRepository) GetBySource(sourceType string, sourceID int64) ([]models.Task, error) {
	var tasks []models.Task
	err := r.QuerySlice(database.SelectTasksBySource, []any{sourceType, sourceID}, func(scan func(dest ...any) error) error {
		t, err := models.ScanTask(scan)
		if err != nil {
			return err
		}
		tasks = append(tasks, *t)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}

// Create inserts a new task
func (r *TaskRepository) Create(projectID int64, sourceType string, sourceID *int64, title, description, priority, dueDate string) (int64, error) {
	result, err := r.Exec(
		database.InsertTask,
		projectID, sourceType, sourceID, title, description, priority, dueDate,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update updates an existing task
func (r *TaskRepository) Update(id int64, title, description, status, priority, dueDate string) error {
	result, err := r.Exec(
		database.UpdateTask,
		title, description, status, priority, dueDate, id,
	)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("task")
	}
	
	return nil
}

// UpdateStatus updates only the status of a task
func (r *TaskRepository) UpdateStatus(id int64, status string) error {
	result, err := r.Exec(database.UpdateTaskStatus, status, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("task")
	}
	
	return nil
}

// Delete removes a task
func (r *TaskRepository) Delete(id int64) error {
	result, err := r.Exec(database.DeleteTask, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("task")
	}
	
	return nil
}

// DeleteBySource removes tasks linked to a specific source
func (r *TaskRepository) DeleteBySource(sourceType string, sourceID int64) error {
	_, err := r.Exec(database.DeleteTasksBySource, sourceType, sourceID)
	return err
}

// Search searches tasks by title or description
func (r *TaskRepository) Search(query string) ([]models.TaskWithContext, error) {
	var tasks []models.TaskWithContext
	err := r.QuerySlice(database.SearchTasks, []any{query, query}, func(scan func(dest ...any) error) error {
		t, err := models.ScanTaskWithContext(scan)
		if err != nil {
			return err
		}
		tasks = append(tasks, *t)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if tasks == nil {
		tasks = []models.TaskWithContext{}
	}
	return tasks, nil
}
