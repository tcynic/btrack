package project

import (
	"context"
	"fmt"
	"log"
	"time"

	"btrack/internal/models"
)

// Update updates project details and recalculates weekly entries if needed
func (s *Service) Update(ctx context.Context, input models.UpdateProjectInput, autoBackupFn func() error) (*models.ProjectWithStats, error) {
	// Verify project exists
	existing, err := s.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Check if we need to recalculate entries (skip for persistent projects)
	needsRecalc := !existing.IsPersistent && (existing.TotalSoldHours != input.TotalSoldHours ||
		existing.SpecialistHours != input.SpecialistHours ||
		existing.StartDate != input.StartDate ||
		existing.EndDate != input.EndDate)

	// Update project
	_, err = s.store.UpdateProject(input)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Recalculate entries if needed
	if needsRecalc {
		// Get next Monday (start of future weeks)
		nextMonday := getCurrentWeekMonday().AddDate(0, 0, 7)
		nextMondayStr := nextMonday.Format("2006-01-02")

		// Parse the new start and end dates
		newStartDate, err := time.Parse("2006-01-02", input.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start date format: %w", err)
		}

		newEndDate, err := time.Parse("2006-01-02", input.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date format: %w", err)
		}

		// Get current weekly entries
		currentEntries, err := s.store.GetWeeklyEntries(input.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get weekly entries: %w", err)
		}

		// If the new start date is after next Monday, recalculate everything from scratch
		if newStartDate.After(nextMonday) || newStartDate.Equal(nextMonday) {
			// Delete all existing entries
			for _, entry := range currentEntries {
				if err := s.store.DeleteWeeklyEntry(input.ID, entry.ID); err != nil {
					return nil, fmt.Errorf("failed to delete weekly entry: %w", err)
				}
			}

			// Calculate new distribution from scratch
			entries, err := s.CalculateDistribution(models.CreateProjectInput{
				Name:            input.Name,
				TotalSoldHours:  input.TotalSoldHours,
				SpecialistHours: input.SpecialistHours,
				StartDate:       input.StartDate,
				EndDate:         input.EndDate,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to calculate distribution: %w", err)
			}

			// Insert new entries
			for _, entry := range entries {
				_, err := s.store.AddWeeklyEntry(input.ID, entry)
				if err != nil {
					return nil, fmt.Errorf("failed to add weekly entry: %w", err)
				}
			}
		} else {
			// Preserve past/current weeks, only recalculate future weeks

			// Calculate sum of past week planned hours
			pastWeekHours := 0
			for _, entry := range currentEntries {
				if entry.WeekStartDate < nextMondayStr {
					pastWeekHours += entry.PlannedHours
				}
			}

			// Delete only future week entries
			if err := s.store.DeleteFutureWeeklyEntries(input.ID, nextMondayStr); err != nil {
				return nil, fmt.Errorf("failed to delete future weekly entries: %w", err)
			}

			// Calculate remaining hours and weeks
			myHours := input.TotalSoldHours - input.SpecialistHours
			remainingHours := myHours - pastWeekHours

			// Only create future entries if there are remaining hours and the end date is in the future
			if remainingHours > 0 && newEndDate.After(nextMonday) {
				// Calculate weeks between next Monday and end date
				futureWeeks := calculateWeeksBetween(nextMonday, newEndDate)
				if futureWeeks < 1 {
					futureWeeks = 1
				}

				// Distribute remaining hours across future weeks (frontloaded)
				baseHours := remainingHours / futureWeeks
				remainder := remainingHours % futureWeeks

				currentMonday := nextMonday
				for i := 0; i < futureWeeks; i++ {
					hours := baseHours
					if i < remainder {
						hours++ // Frontload extra hours
					}

					entry := models.WeeklyEntry{
						WeekStartDate: currentMonday.Format("2006-01-02"),
						WeekNumber:    0,
						PlannedHours:  hours,
						ActualHours:   0,
					}
					_, err := s.store.AddWeeklyEntry(input.ID, entry)
					if err != nil {
						return nil, fmt.Errorf("failed to add weekly entry: %w", err)
					}

					currentMonday = currentMonday.AddDate(0, 0, 7)
				}
			}
		}
	}

	return s.GetByID(ctx, input.ID)
}

// Delete deletes a project and all its nested data
func (s *Service) Delete(ctx context.Context, id int64, autoBackupFn func() error) error {
	// Check if project is persistent
	project, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if project.IsPersistent {
		return models.Forbidden("cannot delete persistent project")
	}

	// Auto-backup before deletion
	if err := autoBackupFn(); err != nil {
		// Log warning but don't fail the operation
		log.Printf("Warning: auto-backup failed before delete: %v", err)
	}

	return s.store.DeleteProject(id)
}

// Restore is no longer needed with JSON store (no soft delete)
// Keeping stub for API compatibility
func (s *Service) Restore(ctx context.Context, id int64) (*models.ProjectWithStats, error) {
	return nil, fmt.Errorf("restore not supported with JSON store")
}

// PermanentlyDelete is the same as Delete with JSON store (no soft delete)
func (s *Service) PermanentlyDelete(ctx context.Context, id int64, autoBackupFn func() error) error {
	return s.Delete(ctx, id, autoBackupFn)
}

// ToggleActive toggles the is_active status of a project
func (s *Service) ToggleActive(ctx context.Context, id int64, autoBackupFn func() error) (*models.ProjectWithStats, error) {
	project, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.Update(ctx, models.UpdateProjectInput{
		ID:              project.ID,
		Name:            project.Name,
		TotalSoldHours:  project.TotalSoldHours,
		SpecialistHours: project.SpecialistHours,
		StartDate:       project.StartDate,
		EndDate:         project.EndDate,
		IsActive:        !project.IsActive,
	}, autoBackupFn)
}

// Helper function
func getCurrentWeekMonday() time.Time {
	today := time.Now()
	weekday := int(today.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	monday := today.AddDate(0, 0, -(weekday - 1))
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
}
