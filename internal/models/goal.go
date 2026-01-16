package models

import (
	"time"

	"btrack/internal/database"
)

// Goal status constants
const (
	GoalStatusPending    = "pending"
	GoalStatusInProgress = "in_progress"
	GoalStatusCompleted  = "completed"
	GoalStatusCancelled  = "cancelled"
)

// Goal represents a project goal
type Goal struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"projectId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	TargetDate  string    `json:"targetDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GoalWithProject includes the parent project name for search results
type GoalWithProject struct {
	Goal
	ProjectName string `json:"projectName"`
}

// ScanGoal scans a database row into a Goal struct.
// Expected columns: id, project_id, title, description, status, target_date, created_at, updated_at
func ScanGoal(scan func(dest ...any) error) (*Goal, error) {
	var g Goal
	var description, targetDate *string
	var createdAt, updatedAt string

	err := scan(
		&g.ID,
		&g.ProjectID,
		&g.Title,
		&description,
		&g.Status,
		&targetDate,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	g.Description = database.NullableString(description)
	g.TargetDate = database.NullableString(targetDate)
	g.CreatedAt = database.ParseTimestamp(createdAt)
	g.UpdatedAt = database.ParseTimestamp(updatedAt)

	return &g, nil
}

// ScanGoalWithProject scans a database row into a GoalWithProject struct.
// Expected columns: id, project_id, title, description, status, target_date, created_at, updated_at, project_name
func ScanGoalWithProject(scan func(dest ...any) error) (*GoalWithProject, error) {
	var g GoalWithProject
	var description, targetDate *string
	var createdAt, updatedAt string

	err := scan(
		&g.ID,
		&g.ProjectID,
		&g.Title,
		&description,
		&g.Status,
		&targetDate,
		&createdAt,
		&updatedAt,
		&g.ProjectName,
	)
	if err != nil {
		return nil, err
	}

	g.Description = database.NullableString(description)
	g.TargetDate = database.NullableString(targetDate)
	g.CreatedAt = database.ParseTimestamp(createdAt)
	g.UpdatedAt = database.ParseTimestamp(updatedAt)

	return &g, nil
}

// CreateGoalInput is the input for creating a new goal
type CreateGoalInput struct {
	ProjectID   int64  `json:"projectId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TargetDate  string `json:"targetDate"`
}

// UpdateGoalInput is the input for updating a goal
type UpdateGoalInput struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	TargetDate  string `json:"targetDate"`
}

// Validate checks if the CreateGoalInput is valid
func (c *CreateGoalInput) Validate() error {
	if c.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

// Validate checks if the UpdateGoalInput is valid
func (u *UpdateGoalInput) Validate() error {
	if u.Title == "" {
		return ErrTitleRequired
	}
	if u.Status != "" && u.Status != GoalStatusPending && u.Status != GoalStatusInProgress && u.Status != GoalStatusCompleted && u.Status != GoalStatusCancelled {
		return ErrInvalidStatus
	}
	return nil
}

// IsValidStatus checks if a status string is valid
func IsValidStatus(status string) bool {
	return status == GoalStatusPending || status == GoalStatusInProgress || status == GoalStatusCompleted || status == GoalStatusCancelled
}
