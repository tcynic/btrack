package store

import (
	"time"
)

// DashboardWeekData represents aggregated data for a single week
type DashboardWeekData struct {
	WeekStartDate     string `json:"weekStartDate"`
	TotalPlannedHours int    `json:"totalPlannedHours"`
	TotalActualHours  int    `json:"totalActualHours"`
	ProjectCount      int    `json:"projectCount"`
}

// MonthlyTrend represents aggregated data for a single month
type MonthlyTrend struct {
	Month        string `json:"month"` // Format: "2024-01"
	PlannedHours int    `json:"plannedHours"`
	ActualHours  int    `json:"actualHours"`
	ProjectCount int    `json:"projectCount"`
	Variance     int    `json:"variance"` // actual - planned
}

// VarianceReport represents planned vs actual for a single project
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

// GetWeeklyAggregates returns aggregated data for all active projects by week
func (s *Store) GetWeeklyAggregates(startDate, endDate string) ([]DashboardWeekData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Map to group entries by week
	weekMap := make(map[string]*DashboardWeekData)

	for _, project := range s.data.Projects {
		if !project.IsActive {
			continue
		}

		for _, entry := range project.WeeklyEntries {
			if entry.WeekStartDate >= startDate && entry.WeekStartDate <= endDate {
				data, exists := weekMap[entry.WeekStartDate]
				if !exists {
					data = &DashboardWeekData{
						WeekStartDate: entry.WeekStartDate,
						ProjectCount:  0,
					}
					weekMap[entry.WeekStartDate] = data
				}

				data.TotalPlannedHours += entry.PlannedHours
				data.TotalActualHours += entry.ActualHours
				data.ProjectCount++
			}
		}
	}

	// Convert map to sorted slice
	var result []DashboardWeekData
	for _, data := range weekMap {
		result = append(result, *data)
	}

	// Sort by week start date
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].WeekStartDate > result[j].WeekStartDate {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if result == nil {
		result = []DashboardWeekData{}
	}

	return result, nil
}

// GetMonthlyTrends returns aggregated data by month
func (s *Store) GetMonthlyTrends(monthsBack int) ([]MonthlyTrend, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Calculate start date
	now := time.Now()
	startDate := now.AddDate(0, -monthsBack, 0)
	startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	startDateStr := startDate.Format("2006-01-02")

	// Map to group by month
	monthMap := make(map[string]*MonthlyTrend)
	projectsPerMonth := make(map[string]map[int64]bool) // track unique projects per month

	for _, project := range s.data.Projects {
		if !project.IsActive {
			continue
		}

		for _, entry := range project.WeeklyEntries {
			if entry.WeekStartDate >= startDateStr {
				// Extract month from week start date (YYYY-MM)
				month := entry.WeekStartDate[:7]

				data, exists := monthMap[month]
				if !exists {
					data = &MonthlyTrend{
						Month: month,
					}
					monthMap[month] = data
					projectsPerMonth[month] = make(map[int64]bool)
				}

				data.PlannedHours += entry.PlannedHours
				data.ActualHours += entry.ActualHours
				projectsPerMonth[month][project.ID] = true
			}
		}
	}

	// Calculate variance and project counts
	for month, data := range monthMap {
		data.Variance = data.ActualHours - data.PlannedHours
		data.ProjectCount = len(projectsPerMonth[month])
	}

	// Convert map to sorted slice
	var result []MonthlyTrend
	for _, data := range monthMap {
		result = append(result, *data)
	}

	// Sort by month
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Month > result[j].Month {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if result == nil {
		result = []MonthlyTrend{}
	}

	return result, nil
}

// GetVarianceReport returns planned vs actual for projects in a date range
func (s *Store) GetVarianceReport(startDate, endDate string) ([]VarianceReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []VarianceReport

	for _, project := range s.data.Projects {
		if !project.IsActive {
			continue
		}

		planned := 0
		actual := 0

		for _, entry := range project.WeeklyEntries {
			if entry.WeekStartDate >= startDate && entry.WeekStartDate <= endDate {
				planned += entry.PlannedHours
				actual += entry.ActualHours
			}
		}

		// Only include projects with planned hours
		if planned > 0 {
			percentage := 0
			if planned > 0 {
				percentage = int(float64(actual) / float64(planned) * 100)
			}

			result = append(result, VarianceReport{
				ProjectID:   project.ID,
				ProjectName: project.Name,
				Planned:     planned,
				Actual:      actual,
				Variance:    actual - planned,
				Percentage:  percentage,
			})
		}
	}

	// Sort by planned hours (descending)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Planned < result[j].Planned {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if result == nil {
		result = []VarianceReport{}
	}

	return result, nil
}

// GetCapacityUtilization returns weekly capacity utilization data
func (s *Store) GetCapacityUtilization(startDate, endDate string) ([]CapacityWeek, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Map to group by week
	weekMap := make(map[string]*CapacityWeek)
	projectsPerWeek := make(map[string]map[int64]bool)

	for _, project := range s.data.Projects {
		if !project.IsActive {
			continue
		}

		for _, entry := range project.WeeklyEntries {
			if entry.WeekStartDate >= startDate && entry.WeekStartDate <= endDate {
				data, exists := weekMap[entry.WeekStartDate]
				if !exists {
					data = &CapacityWeek{
						WeekStart: entry.WeekStartDate,
					}
					weekMap[entry.WeekStartDate] = data
					projectsPerWeek[entry.WeekStartDate] = make(map[int64]bool)
				}

				data.PlannedHours += entry.PlannedHours
				data.ActualHours += entry.ActualHours
				projectsPerWeek[entry.WeekStartDate][project.ID] = true
			}
		}
	}

	// Calculate utilization and project counts
	for week, data := range weekMap {
		if data.PlannedHours > 0 {
			data.Utilization = float64(data.ActualHours) / float64(data.PlannedHours) * 100
		}
		data.ProjectCount = len(projectsPerWeek[week])
	}

	// Convert map to sorted slice
	var result []CapacityWeek
	for _, data := range weekMap {
		result = append(result, *data)
	}

	// Sort by week start
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].WeekStart > result[j].WeekStart {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if result == nil {
		result = []CapacityWeek{}
	}

	return result, nil
}

// GetGoalStats returns goal statistics across all active projects
func (s *Store) GetGoalStats() (totalGoals, completedGoals int, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, project := range s.data.Projects {
		if !project.IsActive {
			continue
		}

		for _, goal := range project.Goals {
			totalGoals++
			if goal.Status == "completed" {
				completedGoals++
			}
		}
	}

	return totalGoals, completedGoals, nil
}

// GetGoalStatsByProject returns goal statistics for a specific project
func (s *Store) GetGoalStatsByProject(projectID int64) (map[string]int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, nil // Return empty stats instead of error
	}

	stats := make(map[string]int)

	for _, goal := range s.data.Projects[projectIdx].Goals {
		stats[goal.Status]++
	}

	return stats, nil
}

// GetActiveProjectCount returns the number of active projects
func (s *Store) GetActiveProjectCount() (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, project := range s.data.Projects {
		if project.IsActive {
			count++
		}
	}

	return count, nil
}
