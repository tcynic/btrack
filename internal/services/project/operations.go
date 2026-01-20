package project

import (
	"context"
	"fmt"
	"log"
	"time"

	"btrack/internal/database"
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

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Update project
	isActive := 0
	if input.IsActive {
		isActive = 1
	}

	_, err = tx.Exec(database.UpdateProject,
		input.Name,
		input.TotalSoldHours,
		input.SpecialistHours,
		input.StartDate,
		input.EndDate,
		isActive,
		input.ID,
	)
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

		// If the new start date is after next Monday, we need to delete everything and recalculate
		if newStartDate.After(nextMonday) || newStartDate.Equal(nextMonday) {
			// Delete all existing entries
			_, err = tx.Exec(database.DeleteWeeklyEntriesByProject, input.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to delete weekly entries: %w", err)
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
				_, err := tx.Exec(database.InsertWeeklyEntry,
					input.ID,
					entry.WeekStartDate,
					entry.WeekNumber,
					entry.PlannedHours,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to insert weekly entry: %w", err)
				}
			}
		} else {
			// Preserve past/current weeks, only recalculate future weeks

			// Get sum of planned hours for past/current weeks
			var pastWeekHours int
			err = tx.QueryRow(database.SelectPastWeeklyHours, input.ID, nextMondayStr).Scan(&pastWeekHours)
			if err != nil {
				return nil, fmt.Errorf("failed to get past week hours: %w", err)
			}

			// Delete only future week entries
			_, err = tx.Exec(database.DeleteFutureWeeklyEntries, input.ID, nextMondayStr)
			if err != nil {
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

					_, err := tx.Exec(database.InsertWeeklyEntry,
						input.ID,
						currentMonday.Format("2006-01-02"),
						0, // Week number will be recalculated; set to 0 for now
						hours,
					)
					if err != nil {
						return nil, fmt.Errorf("failed to insert weekly entry: %w", err)
					}

					currentMonday = currentMonday.AddDate(0, 0, 7)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return s.GetByID(ctx, input.ID)
}

// Delete soft-deletes a project (sets deleted_at timestamp)
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

	result, err := s.db.Exec(database.SoftDeleteProject, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.NotFound("project")
	}

	return nil
}

// Restore restores a soft-deleted project
func (s *Service) Restore(ctx context.Context, id int64) (*models.ProjectWithStats, error) {
	result, err := s.db.Exec(database.RestoreProject, id)
	if err != nil {
		return nil, fmt.Errorf("failed to restore project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, models.NotFound("project")
	}

	// Return the restored project
	// Note: We need to query without the deleted_at filter
	row := s.db.QueryRow(`
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, is_persistent, created_at, updated_at
		FROM projects WHERE id = ?
	`, id)
	p, err := models.ScanProject(row.Scan)
	if err != nil {
		return nil, models.NotFound("project")
	}

	stats, err := s.getStats(p.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project stats: %w", err)
	}

	return &models.ProjectWithStats{
		Project:           *p,
		MyHours:           p.TotalSoldHours - p.SpecialistHours,
		TotalWeeks:        stats.TotalWeeks,
		TotalPlannedHours: stats.TotalPlanned,
		TotalActualHours:  stats.TotalActual,
	}, nil
}

// PermanentlyDelete permanently removes a project and its weekly entries
func (s *Service) PermanentlyDelete(ctx context.Context, id int64, autoBackupFn func() error) error {
	// Check if project is persistent
	// Query without deleted_at filter to check persistence
	var isPersistent int
	err := s.db.QueryRow(`
		SELECT is_persistent FROM projects WHERE id = ?
	`, id).Scan(&isPersistent)
	if err != nil {
		return models.NotFound("project")
	}

	if isPersistent == 1 {
		return models.Forbidden("cannot permanently delete persistent project")
	}

	// Auto-backup before permanent deletion
	if err := autoBackupFn(); err != nil {
		// Log warning but don't fail the operation
		log.Printf("Warning: auto-backup failed before permanent delete: %v", err)
	}

	result, err := s.db.Exec(database.PermanentlyDeleteProject, id)
	if err != nil {
		return fmt.Errorf("failed to permanently delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.NotFound("project")
	}

	return nil
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
