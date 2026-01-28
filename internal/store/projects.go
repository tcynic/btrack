package store

import (
	"fmt"
	"strings"
	"time"

	"btrack/internal/models"
)

// CreateProject creates a new project with empty nested collections
func (s *Store) CreateProject(input models.CreateProjectInput) (*models.Project, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := input.Validate(); err != nil {
		return nil, err
	}

	now := time.Now().Format(time.RFC3339)
	project := ProjectWithNested{
		Project: models.Project{
			ID:              s.data.NextIDs["project"],
			Name:            input.Name,
			TotalSoldHours:  input.TotalSoldHours,
			SpecialistHours: input.SpecialistHours,
			StartDate:       input.StartDate,
			EndDate:         input.EndDate,
			IsActive:        true,
			IsPersistent:    false,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		WeeklyEntries: []models.WeeklyEntry{},
		Meetings:      []models.Meeting{},
		Notes:         []models.Note{},
		Goals:         []models.Goal{},
		Tasks:         []models.Task{},
	}

	s.data.NextIDs["project"]++
	s.data.Projects = append(s.data.Projects, project)

	if err := s.saveUnlocked(); err != nil {
		return nil, err
	}

	return &project.Project, nil
}

// GetProject retrieves a project by ID
func (s *Store) GetProject(id int64) (*ProjectWithNested, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.data.Projects {
		if s.data.Projects[i].ID == id {
			return &s.data.Projects[i], nil
		}
	}

	return nil, models.NotFound("project")
}

// GetAllProjects returns all projects, optionally filtered by active status
func (s *Store) GetAllProjects(activeOnly bool) ([]ProjectWithNested, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !activeOnly {
		// Return copy of all projects
		result := make([]ProjectWithNested, len(s.data.Projects))
		copy(result, s.data.Projects)
		return result, nil
	}

	// Filter for active only
	var result []ProjectWithNested
	for _, p := range s.data.Projects {
		if p.IsActive {
			result = append(result, p)
		}
	}

	if result == nil {
		result = []ProjectWithNested{}
	}

	return result, nil
}

// UpdateProject updates an existing project
func (s *Store) UpdateProject(input models.UpdateProjectInput) (*models.Project, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find project
	var projectIdx = -1
	for i := range s.data.Projects {
		if s.data.Projects[i].ID == input.ID {
			projectIdx = i
			break
		}
	}

	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	// Update fields
	now := time.Now().Format(time.RFC3339)
	s.data.Projects[projectIdx].Name = input.Name
	s.data.Projects[projectIdx].TotalSoldHours = input.TotalSoldHours
	s.data.Projects[projectIdx].SpecialistHours = input.SpecialistHours
	s.data.Projects[projectIdx].StartDate = input.StartDate
	s.data.Projects[projectIdx].EndDate = input.EndDate
	s.data.Projects[projectIdx].IsActive = input.IsActive
	s.data.Projects[projectIdx].UpdatedAt = now

	if err := s.saveUnlocked(); err != nil {
		return nil, err
	}

	return &s.data.Projects[projectIdx].Project, nil
}

// DeleteProject removes a project and all its nested data
func (s *Store) DeleteProject(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find project index
	var projectIdx = -1
	for i := range s.data.Projects {
		if s.data.Projects[i].ID == id {
			projectIdx = i
			break
		}
	}

	if projectIdx == -1 {
		return models.NotFound("project")
	}

	// Remove project (and all nested data automatically)
	s.data.Projects = append(s.data.Projects[:projectIdx], s.data.Projects[projectIdx+1:]...)

	return s.saveUnlocked()
}

// SearchProjects searches for projects by name
func (s *Store) SearchProjects(query string) ([]models.Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(query)
	var result []models.Project

	for _, p := range s.data.Projects {
		if strings.Contains(strings.ToLower(p.Name), query) {
			result = append(result, p.Project)
		}
	}

	if result == nil {
		result = []models.Project{}
	}

	return result, nil
}

// GetProjectStats computes statistics for a project from its weekly entries
func (s *Store) GetProjectStats(projectID int64) (totalPlanned, totalActual, totalWeeks int, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Find project
	var project *ProjectWithNested
	for i := range s.data.Projects {
		if s.data.Projects[i].ID == projectID {
			project = &s.data.Projects[i]
			break
		}
	}

	if project == nil {
		return 0, 0, 0, models.NotFound("project")
	}

	// Sum up weekly entries
	for _, entry := range project.WeeklyEntries {
		totalPlanned += entry.PlannedHours
		totalActual += entry.ActualHours
	}
	totalWeeks = len(project.WeeklyEntries)

	return totalPlanned, totalActual, totalWeeks, nil
}

// findProjectIndex returns the index of a project, or -1 if not found
func (s *Store) findProjectIndex(id int64) int {
	for i := range s.data.Projects {
		if s.data.Projects[i].ID == id {
			return i
		}
	}
	return -1
}

// CreatePersistentProject creates a persistent project (used for seeding)
func (s *Store) CreatePersistentProject(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if already exists
	for _, p := range s.data.Projects {
		if p.Name == name && p.IsPersistent {
			return nil // Already exists
		}
	}

	now := time.Now().Format(time.RFC3339)
	project := ProjectWithNested{
		Project: models.Project{
			ID:              s.data.NextIDs["project"],
			Name:            name,
			TotalSoldHours:  0,
			SpecialistHours: 0,
			StartDate:       "1900-01-01",
			EndDate:         "2099-12-31",
			IsActive:        true,
			IsPersistent:    true,
			CreatedAt:       now,
			UpdatedAt:       now,
		},
		WeeklyEntries: []models.WeeklyEntry{},
		Meetings:      []models.Meeting{},
		Notes:         []models.Note{},
		Goals:         []models.Goal{},
		Tasks:         []models.Task{},
	}

	s.data.NextIDs["project"]++
	s.data.Projects = append(s.data.Projects, project)

	return s.saveUnlocked()
}

// SeedPersistentProjects creates the default persistent projects
func (s *Store) SeedPersistentProjects() error {
	persistentProjects := []string{"Management", "Internal Projects"}
	
	for _, name := range persistentProjects {
		if err := s.CreatePersistentProject(name); err != nil {
			return fmt.Errorf("failed to seed persistent project %s: %w", name, err)
		}
	}
	
	return nil
}
