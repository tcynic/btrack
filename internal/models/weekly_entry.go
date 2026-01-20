package models

// WeeklyEntry represents planned and actual hours for a specific week
type WeeklyEntry struct {
	ID            int64  `json:"id"`
	ProjectID     int64  `json:"projectId"`
	WeekStartDate string `json:"weekStartDate"`
	WeekNumber    int    `json:"weekNumber"`
	PlannedHours  int    `json:"plannedHours"`
	ActualHours   int    `json:"actualHours"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// WeeklyEntryWithStatus includes computed status information
type WeeklyEntryWithStatus struct {
	WeeklyEntry
	Variance   int    `json:"variance"`
	Status     string `json:"status"`
	IsPastWeek bool   `json:"isPastWeek"`
}

// WeeklyEntryWithProject includes project name for cross-project queries
type WeeklyEntryWithProject struct {
	WeeklyEntryWithStatus
	ProjectName string `json:"projectName"`
}

// UpdateActualHoursInput is the input for updating actual hours
type UpdateActualHoursInput struct {
	EntryID     int64 `json:"entryId"`
	ActualHours int   `json:"actualHours"`
}

// ScanWeeklyEntry scans a database row into a WeeklyEntry struct.
// Expected columns: id, project_id, week_start_date, week_number, planned_hours, actual_hours, created_at, updated_at
func ScanWeeklyEntry(scan func(dest ...any) error) (*WeeklyEntry, error) {
	var e WeeklyEntry
	var createdAt, updatedAt string

	err := scan(
		&e.ID,
		&e.ProjectID,
		&e.WeekStartDate,
		&e.WeekNumber,
		&e.PlannedHours,
		&e.ActualHours,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	e.CreatedAt = createdAt
	e.UpdatedAt = updatedAt

	return &e, nil
}

// ScanWeeklyEntryWithProject scans a database row into a WeeklyEntryWithProject struct.
// Expected columns: id, project_id, week_start_date, week_number, planned_hours, actual_hours, created_at, updated_at, project_name
func ScanWeeklyEntryWithProject(scan func(dest ...any) error, isPastWeek bool) (*WeeklyEntryWithProject, error) {
	var e WeeklyEntryWithProject
	var createdAt, updatedAt string

	err := scan(
		&e.ID,
		&e.ProjectID,
		&e.WeekStartDate,
		&e.WeekNumber,
		&e.PlannedHours,
		&e.ActualHours,
		&createdAt,
		&updatedAt,
		&e.ProjectName,
	)
	if err != nil {
		return nil, err
	}

	e.CreatedAt = createdAt
	e.UpdatedAt = updatedAt
	e.IsPastWeek = isPastWeek
	e.CalculateStatus()

	return &e, nil
}

// CalculateStatus computes the status based on planned vs actual hours
func (w *WeeklyEntryWithStatus) CalculateStatus() {
	w.Variance = w.PlannedHours - w.ActualHours

	if w.ActualHours == 0 && !w.IsPastWeek {
		w.Status = "pending"
	} else if w.Variance > 0 {
		w.Status = "under"
	} else if w.Variance < 0 {
		w.Status = "over"
	} else {
		w.Status = "on-track"
	}
}
