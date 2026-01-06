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
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ProjectWithStats includes computed statistics
type ProjectWithStats struct {
	Project
	MyHours           int `json:"myHours"`
	TotalWeeks        int `json:"totalWeeks"`
	TotalPlannedHours int `json:"totalPlannedHours"`
	TotalActualHours  int `json:"totalActualHours"`
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
