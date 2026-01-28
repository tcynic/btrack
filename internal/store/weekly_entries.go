package store

import (
	"time"

	"btrack/internal/models"
)

// AddWeeklyEntry adds a weekly entry to a project
func (s *Store) AddWeeklyEntry(projectID int64, entry models.WeeklyEntry) (*models.WeeklyEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	now := time.Now().Format(time.RFC3339)
	entry.ID = s.data.NextIDs["weekly_entry"]
	entry.ProjectID = projectID
	entry.CreatedAt = now
	entry.UpdatedAt = now

	s.data.NextIDs["weekly_entry"]++
	s.data.Projects[projectIdx].WeeklyEntries = append(s.data.Projects[projectIdx].WeeklyEntries, entry)

	if err := s.saveUnlocked(); err != nil {
		return nil, err
	}

	return &entry, nil
}

// GetWeeklyEntries returns all weekly entries for a project
func (s *Store) GetWeeklyEntries(projectID int64) ([]models.WeeklyEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	entries := s.data.Projects[projectIdx].WeeklyEntries
	if entries == nil {
		entries = []models.WeeklyEntry{}
	}

	return entries, nil
}

// GetWeeklyEntry returns a specific weekly entry
func (s *Store) GetWeeklyEntry(projectID, entryID int64) (*models.WeeklyEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].WeeklyEntries {
		if s.data.Projects[projectIdx].WeeklyEntries[i].ID == entryID {
			return &s.data.Projects[projectIdx].WeeklyEntries[i], nil
		}
	}

	return nil, models.NotFound("weekly entry")
}

// UpdateWeeklyEntry updates a weekly entry
func (s *Store) UpdateWeeklyEntry(projectID, entryID int64, plannedHours, actualHours int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].WeeklyEntries {
		if s.data.Projects[projectIdx].WeeklyEntries[i].ID == entryID {
			s.data.Projects[projectIdx].WeeklyEntries[i].PlannedHours = plannedHours
			s.data.Projects[projectIdx].WeeklyEntries[i].ActualHours = actualHours
			s.data.Projects[projectIdx].WeeklyEntries[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return s.saveUnlocked()
		}
	}

	return models.NotFound("weekly entry")
}

// UpdateActualHours updates only the actual hours for a weekly entry
func (s *Store) UpdateActualHours(projectID, entryID int64, actualHours int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].WeeklyEntries {
		if s.data.Projects[projectIdx].WeeklyEntries[i].ID == entryID {
			s.data.Projects[projectIdx].WeeklyEntries[i].ActualHours = actualHours
			s.data.Projects[projectIdx].WeeklyEntries[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return s.saveUnlocked()
		}
	}

	return models.NotFound("weekly entry")
}

// DeleteWeeklyEntry removes a weekly entry
func (s *Store) DeleteWeeklyEntry(projectID, entryID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].WeeklyEntries {
		if s.data.Projects[projectIdx].WeeklyEntries[i].ID == entryID {
			// Remove entry
			entries := s.data.Projects[projectIdx].WeeklyEntries
			s.data.Projects[projectIdx].WeeklyEntries = append(entries[:i], entries[i+1:]...)
			return s.saveUnlocked()
		}
	}

	return models.NotFound("weekly entry")
}

// DeleteFutureWeeklyEntries removes future weekly entries for a project starting from a date
func (s *Store) DeleteFutureWeeklyEntries(projectID int64, fromDate string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	var kept []models.WeeklyEntry
	for _, entry := range s.data.Projects[projectIdx].WeeklyEntries {
		if entry.WeekStartDate < fromDate {
			kept = append(kept, entry)
		}
	}

	s.data.Projects[projectIdx].WeeklyEntries = kept
	return s.saveUnlocked()
}
