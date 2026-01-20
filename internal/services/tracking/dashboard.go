package tracking

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
	TotalGoals           int `json:"totalGoals"`
	CompletedGoals       int `json:"completedGoals"`
}

// GetDashboardData returns aggregated weekly data for the dashboard
func (s *Service) GetDashboardData(weeksBack, weeksForward int) ([]DashboardWeekData, error) {
	currentMonday := getCurrentWeekMonday()

	// Calculate date range
	startDate := currentMonday.AddDate(0, 0, -7*weeksBack)
	endDate := currentMonday.AddDate(0, 0, 7*weeksForward)

	rows, err := s.db.Query(database.SelectWeeklyAggregates,
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
func (s *Service) GetDashboardSummary(activeProjects []interface{}) (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	// Get active project count
	err := s.db.QueryRow(database.SelectActiveProjectCount).Scan(&summary.TotalActiveProjects)
	if err != nil {
		return nil, fmt.Errorf("failed to get active project count: %w", err)
	}

	// Count at-risk projects from provided projects
	// activeProjects is expected to be []models.ProjectWithStats but passed as interface{} for flexibility
	for _, p := range activeProjects {
		if proj, ok := p.(models.ProjectWithStats); ok {
			if proj.Health.Status == models.HealthAtRisk || proj.Health.Status == models.HealthOverBudget {
				summary.AtRiskProjects++
			}
		}
	}

	currentMonday := getCurrentWeekMonday()
	nextMonday := currentMonday.AddDate(0, 0, 7)

	// Get this week's data
	thisWeekData, err := s.getWeekData(currentMonday.Format("2006-01-02"))
	if err == nil && thisWeekData != nil {
		summary.TotalPlannedThisWeek = thisWeekData.TotalPlannedHours
		summary.TotalActualThisWeek = thisWeekData.TotalActualHours
	}

	// Get next week's data
	nextWeekData, err := s.getWeekData(nextMonday.Format("2006-01-02"))
	if err == nil && nextWeekData != nil {
		summary.TotalPlannedNextWeek = nextWeekData.TotalPlannedHours
	}

	// Get goal metrics across all active projects
	err = s.db.QueryRow(database.SelectAllGoalStats).Scan(
		&summary.TotalGoals,
		&summary.CompletedGoals,
	)
	// Ignore errors for goal stats, just keep as 0

	return summary, nil
}

// getWeekData retrieves aggregated data for a specific week
func (s *Service) getWeekData(weekStartDate string) (*DashboardWeekData, error) {
	// Use a small range query for a single week
	endDate, _ := time.Parse("2006-01-02", weekStartDate)
	endDateStr := endDate.AddDate(0, 0, 6).Format("2006-01-02")

	rows, err := s.db.Query(database.SelectWeeklyAggregates, weekStartDate, endDateStr)
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
