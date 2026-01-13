package models

import "time"

// WeeklyEntry represents planned and actual hours for a specific week
type WeeklyEntry struct {
	ID            int64     `json:"id"`
	ProjectID     int64     `json:"projectId"`
	WeekStartDate string    `json:"weekStartDate"`
	WeekNumber    int       `json:"weekNumber"`
	PlannedHours  int       `json:"plannedHours"`
	ActualHours   int       `json:"actualHours"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
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
