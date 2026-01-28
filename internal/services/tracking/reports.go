package tracking

import (
	"time"

	"btrack/internal/store"
)

// Type aliases for store types
type MonthlyTrend = store.MonthlyTrend
type VarianceReport = store.VarianceReport
type CapacityWeek = store.CapacityWeek

// GetMonthlyTrends returns aggregated data by month for the specified number of months back
func (s *Service) GetMonthlyTrends(monthsBack int) ([]MonthlyTrend, error) {
	if monthsBack <= 0 {
		monthsBack = 6 // Default to 6 months
	}
	return s.store.GetMonthlyTrends(monthsBack)
}

// GetVarianceReport returns planned vs actual analysis for projects in a date range
func (s *Service) GetVarianceReport(startDate, endDate string) ([]VarianceReport, error) {
	return s.store.GetVarianceReport(startDate, endDate)
}

// GetCapacityUtilization returns weekly capacity utilization data for a date range
func (s *Service) GetCapacityUtilization(startDate, endDate string) ([]CapacityWeek, error) {
	return s.store.GetCapacityUtilization(startDate, endDate)
}
