package main

import (
	"fmt"
	"time"

	"btrack/internal/database"
)

// MonthlyTrend represents aggregated data for a single month
type MonthlyTrend struct {
	Month        string `json:"month"` // Format: "2024-01"
	PlannedHours int    `json:"plannedHours"`
	ActualHours  int    `json:"actualHours"`
	ProjectCount int    `json:"projectCount"`
	Variance     int    `json:"variance"` // actual - planned
}

// VarianceReport represents planned vs actual for a single project in a date range
type VarianceReport struct {
	ProjectID   int64  `json:"projectId"`
	ProjectName string `json:"projectName"`
	Planned     int    `json:"planned"`
	Actual      int    `json:"actual"`
	Variance    int    `json:"variance"`
	Percentage  int    `json:"percentage"` // % of planned hours used
}

// CapacityWeek represents weekly capacity utilization
type CapacityWeek struct {
	WeekStart    string  `json:"weekStart"`
	PlannedHours int     `json:"plannedHours"`
	ActualHours  int     `json:"actualHours"`
	Utilization  float64 `json:"utilization"` // percentage
	ProjectCount int     `json:"projectCount"`
}

// GetMonthlyTrends returns aggregated data by month for the specified number of months back
func (a *App) GetMonthlyTrends(monthsBack int) ([]MonthlyTrend, error) {
	if monthsBack <= 0 {
		monthsBack = 6 // Default to 6 months
	}

	// Calculate start date (first day of month, N months ago)
	now := time.Now()
	startDate := now.AddDate(0, -monthsBack, 0)
	startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	rows, err := a.db.Query(database.SelectMonthlyTrends, startDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly trends: %w", err)
	}
	defer rows.Close()

	var trends []MonthlyTrend
	for rows.Next() {
		var t MonthlyTrend
		if err := rows.Scan(&t.Month, &t.PlannedHours, &t.ActualHours, &t.ProjectCount); err != nil {
			return nil, fmt.Errorf("failed to scan monthly trend: %w", err)
		}
		t.Variance = t.ActualHours - t.PlannedHours
		trends = append(trends, t)
	}

	if trends == nil {
		trends = []MonthlyTrend{}
	}

	return trends, nil
}

// GetVarianceReport returns planned vs actual analysis for projects in a date range
func (a *App) GetVarianceReport(startDate, endDate string) ([]VarianceReport, error) {
	// Validate date formats
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	rows, err := a.db.Query(database.SelectVarianceReport, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query variance report: %w", err)
	}
	defer rows.Close()

	var reports []VarianceReport
	for rows.Next() {
		var r VarianceReport
		if err := rows.Scan(&r.ProjectID, &r.ProjectName, &r.Planned, &r.Actual); err != nil {
			return nil, fmt.Errorf("failed to scan variance report: %w", err)
		}
		r.Variance = r.Actual - r.Planned
		if r.Planned > 0 {
			r.Percentage = int(float64(r.Actual) / float64(r.Planned) * 100)
		}
		reports = append(reports, r)
	}

	if reports == nil {
		reports = []VarianceReport{}
	}

	return reports, nil
}

// GetCapacityUtilization returns weekly capacity utilization data for a date range
func (a *App) GetCapacityUtilization(startDate, endDate string) ([]CapacityWeek, error) {
	// Validate date formats
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	rows, err := a.db.Query(database.SelectCapacityUtilization, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query capacity utilization: %w", err)
	}
	defer rows.Close()

	var weeks []CapacityWeek
	for rows.Next() {
		var w CapacityWeek
		if err := rows.Scan(&w.WeekStart, &w.PlannedHours, &w.ActualHours, &w.ProjectCount); err != nil {
			return nil, fmt.Errorf("failed to scan capacity week: %w", err)
		}
		if w.PlannedHours > 0 {
			w.Utilization = float64(w.ActualHours) / float64(w.PlannedHours) * 100
		}
		weeks = append(weeks, w)
	}

	if weeks == nil {
		weeks = []CapacityWeek{}
	}

	return weeks, nil
}
