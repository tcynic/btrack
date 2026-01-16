package models

import (
	"time"

	"btrack/internal/database"
)

// Task status constants
const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusCancelled  = "cancelled"
)

// Task priority constants
const (
	TaskPriorityLow    = "low"
	TaskPriorityMedium = "medium"
	TaskPriorityHigh   = "high"
)

// Task source type constants
const (
	TaskSourceMeeting    = "meeting"
	TaskSourceNote       = "note"
	TaskSourceStandalone = "standalone"
)

// Task represents a task/action item
type Task struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"projectId"`
	SourceType  string    `json:"sourceType"`
	SourceID    *int64    `json:"sourceId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     string    `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TaskWithContext extends Task with project name and source title
type TaskWithContext struct {
	Task
	ProjectName string `json:"projectName"`
	SourceTitle string `json:"sourceTitle"`
}

// ScanTask scans a database row into a Task struct.
// Expected columns: id, project_id, source_type, source_id, title, description, status, priority, due_date, created_at, updated_at
func ScanTask(scan func(dest ...any) error) (*Task, error) {
	var t Task
	var sourceID *int64
	var description, dueDate *string
	var createdAt, updatedAt string

	err := scan(
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
		return nil, err
	}

	t.SourceID = sourceID
	t.Description = database.NullableString(description)
	t.DueDate = database.NullableString(dueDate)
	t.CreatedAt = database.ParseTimestamp(createdAt)
	t.UpdatedAt = database.ParseTimestamp(updatedAt)

	return &t, nil
}

// ScanTaskWithContext scans a database row into a TaskWithContext struct.
// Expected columns: id, project_id, source_type, source_id, title, description, status, priority, due_date, created_at, updated_at, project_name, source_title
func ScanTaskWithContext(scan func(dest ...any) error) (*TaskWithContext, error) {
	var t TaskWithContext
	var sourceID *int64
	var description, dueDate, sourceTitle *string
	var createdAt, updatedAt string

	err := scan(
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
		return nil, err
	}

	t.SourceID = sourceID
	t.Description = database.NullableString(description)
	t.DueDate = database.NullableString(dueDate)
	t.SourceTitle = database.NullableString(sourceTitle)
	t.CreatedAt = database.ParseTimestamp(createdAt)
	t.UpdatedAt = database.ParseTimestamp(updatedAt)

	return &t, nil
}

// CreateTaskInput is the input for creating a new task
type CreateTaskInput struct {
	ProjectID   int64   `json:"projectId"`
	SourceType  string  `json:"sourceType"`
	SourceID    *int64  `json:"sourceId"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	DueDate     string  `json:"dueDate"`
}

// UpdateTaskInput is the input for updating a task
type UpdateTaskInput struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"dueDate"`
}

// Validate checks if the CreateTaskInput is valid
func (c *CreateTaskInput) Validate() error {
	if c.Title == "" {
		return ErrTitleRequired
	}
	if c.SourceType != "" && !IsValidSourceType(c.SourceType) {
		return ErrInvalidSourceType
	}
	if c.Priority != "" && !IsValidPriority(c.Priority) {
		return ErrInvalidPriority
	}
	// Set defaults
	if c.SourceType == "" {
		c.SourceType = TaskSourceStandalone
	}
	if c.Priority == "" {
		c.Priority = TaskPriorityMedium
	}
	return nil
}

// Validate checks if the UpdateTaskInput is valid
func (u *UpdateTaskInput) Validate() error {
	if u.Title == "" {
		return ErrTitleRequired
	}
	if u.Status != "" && !IsValidTaskStatus(u.Status) {
		return ErrInvalidStatus
	}
	if u.Priority != "" && !IsValidPriority(u.Priority) {
		return ErrInvalidPriority
	}
	return nil
}

// IsValidTaskStatus checks if a status string is valid
func IsValidTaskStatus(status string) bool {
	return status == TaskStatusPending || status == TaskStatusInProgress || status == TaskStatusCompleted || status == TaskStatusCancelled
}

// IsValidPriority checks if a priority string is valid
func IsValidPriority(priority string) bool {
	return priority == TaskPriorityLow || priority == TaskPriorityMedium || priority == TaskPriorityHigh
}

// IsValidSourceType checks if a source type string is valid
func IsValidSourceType(sourceType string) bool {
	return sourceType == TaskSourceMeeting || sourceType == TaskSourceNote || sourceType == TaskSourceStandalone
}
