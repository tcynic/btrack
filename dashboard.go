package main

import (
	"fmt"
	"time"

	"btrack/internal/database"
	"btrack/internal/models"
)

// DashboardWeekData represents aggregated data for a single week across all projects
type DashboardWeekData struct {
	WeekStartDate     string `json:"weekStartDate"`
	TotalPlannedHours int    `json:"totalPlannedHours"`
	TotalActualHours  int    `json:"totalActualHours"`
	ProjectCount      int    `json:"projectCount"`
}

// DashboardSummary provides high-level statistics
type DashboardSummary struct {
	TotalActiveProjects  int `json:"totalActiveProjects"`
	AtRiskProjects       int `json:"atRiskProjects"`
	TotalPlannedThisWeek int `json:"totalPlannedThisWeek"`
	TotalActualThisWeek  int `json:"totalActualThisWeek"`
	TotalPlannedNextWeek int `json:"totalPlannedNextWeek"`
}

// GetDashboardData returns aggregated weekly data for the dashboard
func (a *App) GetDashboardData(weeksBack, weeksForward int) ([]DashboardWeekData, error) {
	currentMonday := getCurrentWeekMonday()

	// Calculate date range
	startDate := currentMonday.AddDate(0, 0, -7*weeksBack)
	endDate := currentMonday.AddDate(0, 0, 7*weeksForward)

	rows, err := a.db.Query(database.SelectWeeklyAggregates,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query dashboard data: %w", err)
	}
	defer rows.Close()

	var data []DashboardWeekData
	for rows.Next() {
		var d DashboardWeekData
		err := rows.Scan(
			&d.WeekStartDate,
			&d.TotalPlannedHours,
			&d.TotalActualHours,
			&d.ProjectCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan dashboard row: %w", err)
		}
		data = append(data, d)
	}

	if data == nil {
		data = []DashboardWeekData{}
	}

	return data, nil
}

// GetDashboardSummary returns high-level statistics for the dashboard
func (a *App) GetDashboardSummary() (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	// Get active project count
	err := a.db.QueryRow(database.SelectActiveProjectCount).Scan(&summary.TotalActiveProjects)
	if err != nil {
		return nil, fmt.Errorf("failed to get active project count: %w", err)
	}

	// Get at-risk projects count
	projects, err := a.GetAllProjects(true)
	if err == nil {
		for _, p := range projects {
			if p.Health.Status == models.HealthAtRisk || p.Health.Status == models.HealthOverBudget {
				summary.AtRiskProjects++
			}
		}
	}

	currentMonday := getCurrentWeekMonday()
	nextMonday := currentMonday.AddDate(0, 0, 7)

	// Get this week's data
	thisWeekData, err := a.getWeekData(currentMonday.Format("2006-01-02"))
	if err == nil && thisWeekData != nil {
		summary.TotalPlannedThisWeek = thisWeekData.TotalPlannedHours
		summary.TotalActualThisWeek = thisWeekData.TotalActualHours
	}

	// Get next week's data
	nextWeekData, err := a.getWeekData(nextMonday.Format("2006-01-02"))
	if err == nil && nextWeekData != nil {
		summary.TotalPlannedNextWeek = nextWeekData.TotalPlannedHours
	}

	return summary, nil
}

// getWeekData retrieves aggregated data for a specific week
func (a *App) getWeekData(weekStartDate string) (*DashboardWeekData, error) {
	// Use a small range query for a single week
	endDate, _ := time.Parse("2006-01-02", weekStartDate)
	endDateStr := endDate.AddDate(0, 0, 6).Format("2006-01-02")

	rows, err := a.db.Query(database.SelectWeeklyAggregates, weekStartDate, endDateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var d DashboardWeekData
		err := rows.Scan(
			&d.WeekStartDate,
			&d.TotalPlannedHours,
			&d.TotalActualHours,
			&d.ProjectCount,
		)
		if err != nil {
			return nil, err
		}
		return &d, nil
	}

	return nil, nil
}
