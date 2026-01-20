package repository

import (
	"btrack/internal/database"
	"btrack/internal/models"
)

// WeeklyEntryRepository handles weekly entry database operations
type WeeklyEntryRepository struct {
	*Repository
}

// NewWeeklyEntryRepository creates a new weekly entry repository
func NewWeeklyEntryRepository(base *Repository) *WeeklyEntryRepository {
	return &WeeklyEntryRepository{Repository: base}
}

// GetByID retrieves a single weekly entry by ID
func (r *WeeklyEntryRepository) GetByID(id int64) (*models.WeeklyEntry, error) {
	var entry *models.WeeklyEntry
	err := r.QueryOne(
		`SELECT id, project_id, week_start_date, week_number, planned_hours, actual_hours, created_at, updated_at
		 FROM weekly_entries WHERE id = ?`,
		[]any{id},
		func(scan func(dest ...any) error) error {
			e, err := models.ScanWeeklyEntry(scan)
			if err != nil {
				return err
			}
			entry = e
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// GetByProject retrieves all weekly entries for a project
func (r *WeeklyEntryRepository) GetByProject(projectID int64) ([]models.WeeklyEntry, error) {
	var entries []models.WeeklyEntry
	err := r.QuerySlice(database.SelectWeeklyEntriesByProject, []any{projectID}, func(scan func(dest ...any) error) error {
		e, err := models.ScanWeeklyEntry(scan)
		if err != nil {
			return err
		}
		entries = append(entries, *e)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if entries == nil {
		entries = []models.WeeklyEntry{}
	}
	return entries, nil
}

// Create inserts a new weekly entry
func (r *WeeklyEntryRepository) Create(projectID int64, weekStartDate string, weekNumber, plannedHours int) (int64, error) {
	result, err := r.Exec(
		database.InsertWeeklyEntry,
		projectID, weekStartDate, weekNumber, plannedHours,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// UpdateActualHours updates the actual hours for a weekly entry
func (r *WeeklyEntryRepository) UpdateActualHours(id int64, actualHours int) error {
	result, err := r.Exec(database.UpdateActualHours, actualHours, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("weekly entry")
	}
	
	return nil
}

// DeleteByProject removes all weekly entries for a project
func (r *WeeklyEntryRepository) DeleteByProject(projectID int64) error {
	_, err := r.Exec(database.DeleteWeeklyEntriesByProject, projectID)
	return err
}

// DeleteFutureEntries removes future weekly entries for a project
func (r *WeeklyEntryRepository) DeleteFutureEntries(projectID int64, fromDate string) error {
	_, err := r.Exec(database.DeleteFutureWeeklyEntries, projectID, fromDate)
	return err
}

// GetPastWeekHours returns sum of planned hours for past weeks
func (r *WeeklyEntryRepository) GetPastWeekHours(projectID int64, beforeDate string) (int, error) {
	var hours int
	err := r.DB().QueryRow(database.SelectPastWeeklyHours, projectID, beforeDate).Scan(&hours)
	if err != nil {
		return 0, models.DatabaseError(err)
	}
	return hours, nil
}
