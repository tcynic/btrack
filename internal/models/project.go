package models

import "time"

// Project represents a client project with sold hours
type Project struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	TotalSoldHours  int       `json:"totalSoldHours"`
	SpecialistHours int       `json:"specialistHours"`
	StartDate       string    `json:"startDate"`
	EndDate         string    `json:"endDate"`
	IsActive        bool      `json:"isActive"`
	IsPersistent    bool      `json:"isPersistent"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// HealthStatus represents the health state of a project
type HealthStatus string

const (
	HealthOnTrack    HealthStatus = "on_track"
	HealthAtRisk     HealthStatus = "at_risk"
	HealthOverBudget HealthStatus = "over_budget"
	HealthCompleted  HealthStatus = "completed"
)

// ProjectHealth contains health status information
type ProjectHealth struct {
	Status   HealthStatus `json:"status"`
	Message  string       `json:"message"`
	Severity string       `json:"severity"` // "info", "warning", "error"
}

// ProjectWithStats includes computed statistics
type ProjectWithStats struct {
	Project
	MyHours           int            `json:"myHours"`
	TotalWeeks        int            `json:"totalWeeks"`
	TotalPlannedHours int            `json:"totalPlannedHours"`
	TotalActualHours  int            `json:"totalActualHours"`
	Health            ProjectHealth  `json:"health"`
}

// CreateProjectInput is the input for creating a new project
type CreateProjectInput struct {
	Name            string `json:"name"`
	TotalSoldHours  int    `json:"totalSoldHours"`
	SpecialistHours int    `json:"specialistHours"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
}

// UpdateProjectInput is the input for updating a project
type UpdateProjectInput struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	TotalSoldHours  int    `json:"totalSoldHours"`
	SpecialistHours int    `json:"specialistHours"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	IsActive        bool   `json:"isActive"`
}

// Validate checks if the CreateProjectInput is valid
func (c *CreateProjectInput) Validate() error {
	if c.Name == "" {
		return ErrNameRequired
	}
	if c.TotalSoldHours <= 0 {
		return ErrInvalidHours
	}
	if c.SpecialistHours < 0 {
		return ErrInvalidHours
	}
	if c.SpecialistHours >= c.TotalSoldHours {
		return ErrSpecialistHoursTooHigh
	}
	if c.StartDate == "" || c.EndDate == "" {
		return ErrDatesRequired
	}
	if c.StartDate > c.EndDate {
		return ErrInvalidDateRange
	}
	return nil
}
