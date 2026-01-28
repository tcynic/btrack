package tracking

import (
	"fmt"
	"time"

	"btrack/internal/models"
)

// GetWeeklyEntries returns all weekly entries for a project with status information
func (s *Service) GetWeeklyEntries(projectID int64) ([]models.WeeklyEntryWithStatus, error) {
	project, err := s.store.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	currentMonday := getCurrentWeekMonday()
	var entries []models.WeeklyEntryWithStatus

	for _, e := range project.WeeklyEntries {
		// Determine if this is a past week (can edit actual hours)
		weekStart, _ := time.Parse("2006-01-02", e.WeekStartDate)
		isPastWeek := weekStart.Before(currentMonday)

		entryWithStatus := models.WeeklyEntryWithStatus{
			WeeklyEntry: e,
			IsPastWeek:  isPastWeek,
		}
		entryWithStatus.CalculateStatus()

		entries = append(entries, entryWithStatus)
	}

	return entries, nil
}

// UpdateActualHours updates the actual hours for a specific weekly entry
func (s *Service) UpdateActualHours(input models.UpdateActualHoursInput) (*models.WeeklyEntryWithStatus, error) {
	if input.ActualHours < 0 {
		return nil, models.ValidationError("actualHours", "hours must be a positive number")
	}

	// Find the entry and its project
	var foundEntry *models.WeeklyEntry
	var projectID int64
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, e := range proj.WeeklyEntries {
			if e.ID == input.EntryID {
				foundEntry = &e
				projectID = proj.ID
				break
			}
		}
		if foundEntry != nil {
			break
		}
	}

	if foundEntry == nil {
		return nil, models.NotFound("weekly entry")
	}

	// Check if this is a past or current week (allow editing)
	weekStart, _ := time.Parse("2006-01-02", foundEntry.WeekStartDate)
	currentMonday := getCurrentWeekMonday()
	nextMonday := currentMonday.AddDate(0, 0, 7)

	if weekStart.After(nextMonday) || weekStart.Equal(nextMonday) {
		return nil, models.Forbidden("cannot edit actual hours for future weeks")
	}

	// Update the actual hours
	err := s.store.UpdateActualHours(projectID, input.EntryID, input.ActualHours)
	if err != nil {
		return nil, err
	}

	// Return the updated entry
	return s.getWeeklyEntry(projectID, input.EntryID)
}

// GetWeeklyEntriesByWeek returns all weekly entries for a specific week across all active projects
func (s *Service) GetWeeklyEntriesByWeek(weekStartDate string) ([]models.WeeklyEntryWithProject, error) {
	// First, ensure all persistent projects have entries for this week
	if err := s.ensurePersistentEntries(weekStartDate); err != nil {
		return nil, fmt.Errorf("failed to ensure persistent entries: %w", err)
	}

	projects, err := s.store.GetAllProjects(true) // only active
	if err != nil {
		return nil, err
	}

	currentMonday := getCurrentWeekMonday()
	var entries []models.WeeklyEntryWithProject

	for _, proj := range projects {
		for _, e := range proj.WeeklyEntries {
			if e.WeekStartDate == weekStartDate {
				weekStart, _ := time.Parse("2006-01-02", e.WeekStartDate)
				isPastWeek := weekStart.Before(currentMonday)

				entryWithProject := models.WeeklyEntryWithProject{
					WeeklyEntryWithStatus: models.WeeklyEntryWithStatus{
						WeeklyEntry: e,
						IsPastWeek:  isPastWeek,
					},
					ProjectName: proj.Name,
				}
				entryWithProject.CalculateStatus()
				entries = append(entries, entryWithProject)
			}
		}
	}

	return entries, nil
}

// ensurePersistentEntries creates weekly entries for all persistent projects if they don't exist
func (s *Service) ensurePersistentEntries(weekStartDate string) error {
	projects, err := s.store.GetAllProjects(true) // only active
	if err != nil {
		return err
	}

	for _, proj := range projects {
		if !proj.IsPersistent {
			continue
		}

		// Check if entry already exists for this week
		exists := false
		for _, e := range proj.WeeklyEntries {
			if e.WeekStartDate == weekStartDate {
				exists = true
				break
			}
		}

		if !exists {
			// Create entry with 0 hours
			entry := models.WeeklyEntry{
				ProjectID:     proj.ID,
				WeekStartDate: weekStartDate,
				WeekNumber:    0,
				PlannedHours:  0,
				ActualHours:   0,
			}
			if _, err := s.store.AddWeeklyEntry(proj.ID, entry); err != nil {
				return fmt.Errorf("failed to create persistent entry: %w", err)
			}
		}
	}

	return nil
}

// getWeeklyEntry retrieves a single weekly entry by ID
func (s *Service) getWeeklyEntry(projectID, entryID int64) (*models.WeeklyEntryWithStatus, error) {
	project, err := s.store.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	for _, e := range project.WeeklyEntries {
		if e.ID == entryID {
			currentMonday := getCurrentWeekMonday()
			weekStart, _ := time.Parse("2006-01-02", e.WeekStartDate)
			isPastWeek := weekStart.Before(currentMonday)

			entryWithStatus := &models.WeeklyEntryWithStatus{
				WeeklyEntry: e,
				IsPastWeek:  isPastWeek,
			}
			entryWithStatus.CalculateStatus()
			return entryWithStatus, nil
		}
	}

	return nil, models.NotFound("weekly entry")
}

// getCurrentWeekMonday returns the Monday of the current week
func getCurrentWeekMonday() time.Time {
	today := time.Now()
	weekday := int(today.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	monday := today.AddDate(0, 0, -(weekday - 1))
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
}
