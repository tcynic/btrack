package store

import (
	"btrack/internal/models"
)

// Data represents the entire application state
type Data struct {
	SchemaVersion int                 `json:"schema_version"`
	NextIDs       map[string]int64    `json:"next_ids"`
	Projects      []ProjectWithNested `json:"projects"`
	Templates     []Template          `json:"templates"`
}

// Template represents a saved project template
type Template struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	TotalSoldHours  int    `json:"total_sold_hours"`
	SpecialistHours int    `json:"specialist_hours"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// ProjectWithNested represents a project with all its nested entities
type ProjectWithNested struct {
	models.Project
	WeeklyEntries []models.WeeklyEntry `json:"weekly_entries"`
	Meetings      []models.Meeting     `json:"meetings"`
	Notes         []models.Note        `json:"notes"`
	Goals         []models.Goal        `json:"goals"`
	Tasks         []models.Task        `json:"tasks"`
}

// NewData creates a new empty Data structure with defaults
func NewData() *Data {
	return &Data{
		SchemaVersion: 1,
		NextIDs: map[string]int64{
			"project":       1,
			"weekly_entry":  1,
			"meeting":       1,
			"note":          1,
			"goal":          1,
			"task":          1,
			"template":      1,
		},
		Projects:  []ProjectWithNested{},
		Templates: []Template{},
	}
}
