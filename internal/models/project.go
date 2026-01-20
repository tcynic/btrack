package models

// Project represents a client project with sold hours
type Project struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	TotalSoldHours  int    `json:"totalSoldHours"`
	SpecialistHours int    `json:"specialistHours"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	IsActive        bool   `json:"isActive"`
	IsPersistent    bool   `json:"isPersistent"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
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

// ScanProject scans a database row into a Project struct.
// Expected columns: id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, is_persistent, created_at, updated_at
func ScanProject(scan func(dest ...any) error) (*Project, error) {
	var p Project
	var isActive, isPersistent int
	var createdAt, updatedAt string

	err := scan(
		&p.ID,
		&p.Name,
		&p.TotalSoldHours,
		&p.SpecialistHours,
		&p.StartDate,
		&p.EndDate,
		&isActive,
		&isPersistent,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.IsActive = isActive == 1
	p.IsPersistent = isPersistent == 1
	p.CreatedAt = createdAt
	p.UpdatedAt = updatedAt

	return &p, nil
}

// Validate checks if the CreateProjectInput is valid
func (c *CreateProjectInput) Validate() error {
	if c.Name == "" {
		return ValidationError("name", "project name is required")
	}
	if c.TotalSoldHours <= 0 {
		return ValidationError("totalSoldHours", "hours must be a positive number")
	}
	if c.SpecialistHours < 0 {
		return ValidationError("specialistHours", "hours must be a positive number")
	}
	if c.SpecialistHours >= c.TotalSoldHours {
		return ValidationError("specialistHours", "specialist hours must be less than total sold hours")
	}
	if c.StartDate == "" || c.EndDate == "" {
		return ValidationError("date", "start date and end date are required")
	}
	if c.StartDate > c.EndDate {
		return ValidationError("date", "start date must be before end date")
	}
	return nil
}
