package main

import (
	"btrack/internal/services/tracking"
)

// DashboardWeekData represents aggregated data for a single week across all projects
type DashboardWeekData = tracking.DashboardWeekData

// DashboardSummary provides high-level statistics
type DashboardSummary = tracking.DashboardSummary

// GetDashboardData returns aggregated weekly data for the dashboard
func (a *App) GetDashboardData(weeksBack, weeksForward int) ([]DashboardWeekData, error) {
	return a.trackingService.GetDashboardData(weeksBack, weeksForward)
}

// GetDashboardSummary returns high-level statistics for the dashboard
func (a *App) GetDashboardSummary() (*DashboardSummary, error) {
	// Get active projects for at-risk calculation
	projects, err := a.GetAllProjects(true)
	if err != nil {
		return nil, err
	}

	// Convert to []interface{} for service method
	projectList := make([]interface{}, len(projects))
	for i, p := range projects {
		projectList[i] = p
	}

	return a.trackingService.GetDashboardSummary(projectList)
}
