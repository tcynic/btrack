package tracking

import (
	"time"

	"btrack/internal/models"
	"btrack/internal/store"
)

// DashboardWeekData is an alias for the store type
type DashboardWeekData = store.DashboardWeekData

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

	return s.store.GetWeeklyAggregates(
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
}

// GetDashboardSummary returns high-level statistics for the dashboard
func (s *Service) GetDashboardSummary(activeProjects []interface{}) (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	// Get active project count
	count, err := s.store.GetActiveProjectCount()
	if err != nil {
		return nil, err
	}
	summary.TotalActiveProjects = count

	// Count at-risk projects from provided projects
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
	if err == nil && len(thisWeekData) > 0 {
		summary.TotalPlannedThisWeek = thisWeekData[0].TotalPlannedHours
		summary.TotalActualThisWeek = thisWeekData[0].TotalActualHours
	}

	// Get next week's data
	nextWeekData, err := s.getWeekData(nextMonday.Format("2006-01-02"))
	if err == nil && len(nextWeekData) > 0 {
		summary.TotalPlannedNextWeek = nextWeekData[0].TotalPlannedHours
	}

	// Get goal metrics across all active projects
	totalGoals, completedGoals, _ := s.store.GetGoalStats()
	summary.TotalGoals = totalGoals
	summary.CompletedGoals = completedGoals

	return summary, nil
}

// getWeekData retrieves aggregated data for a specific week
func (s *Service) getWeekData(weekStartDate string) ([]DashboardWeekData, error) {
	// Use a small range query for a single week
	endDate, _ := time.Parse("2006-01-02", weekStartDate)
	endDateStr := endDate.AddDate(0, 0, 6).Format("2006-01-02")

	return s.store.GetWeeklyAggregates(weekStartDate, endDateStr)
}
