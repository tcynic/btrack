package models

import "time"

// Goal status constants
const (
	GoalStatusPending    = "pending"
	GoalStatusInProgress = "in_progress"
	GoalStatusCompleted  = "completed"
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
	if u.Status != "" && u.Status != GoalStatusPending && u.Status != GoalStatusInProgress && u.Status != GoalStatusCompleted {
		return ErrInvalidStatus
	}
	return nil
}

// IsValidStatus checks if a status string is valid
func IsValidStatus(status string) bool {
	return status == GoalStatusPending || status == GoalStatusInProgress || status == GoalStatusCompleted
}
