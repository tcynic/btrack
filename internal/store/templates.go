package store

import (
	"time"
)

// CreateTemplate creates a new project template
func (s *Store) CreateTemplate(name string, totalSoldHours, specialistHours int) (*Template, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Format(time.RFC3339)
	template := Template{
		ID:              s.data.NextIDs["template"],
		Name:            name,
		TotalSoldHours:  totalSoldHours,
		SpecialistHours: specialistHours,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	s.data.NextIDs["template"]++
	s.data.Templates = append(s.data.Templates, template)

	if err := s.saveUnlocked(); err != nil {
		return nil, err
	}

	return &template, nil
}

// GetTemplates returns all templates
func (s *Store) GetTemplates() ([]Template, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	templates := make([]Template, len(s.data.Templates))
	copy(templates, s.data.Templates)

	if templates == nil {
		templates = []Template{}
	}

	return templates, nil
}

// GetTemplate returns a specific template by ID
func (s *Store) GetTemplate(id int64) (*Template, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.data.Templates {
		if s.data.Templates[i].ID == id {
			return &s.data.Templates[i], nil
		}
	}

	return nil, nil // Return nil instead of error for template not found
}

// DeleteTemplate removes a template
func (s *Store) DeleteTemplate(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data.Templates {
		if s.data.Templates[i].ID == id {
			s.data.Templates = append(s.data.Templates[:i], s.data.Templates[i+1:]...)
			return s.saveUnlocked()
		}
	}

	return nil // Template not found, but not an error
}
