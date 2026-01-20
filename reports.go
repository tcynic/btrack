package main

import (
	"btrack/internal/services/tracking"
)

// MonthlyTrend represents aggregated data for a single month
type MonthlyTrend = tracking.MonthlyTrend

// VarianceReport represents planned vs actual for a single project in a date range
type VarianceReport = tracking.VarianceReport

// CapacityWeek represents weekly capacity utilization
type CapacityWeek = tracking.CapacityWeek

// GetMonthlyTrends returns aggregated data by month for the specified number of months back
func (a *App) GetMonthlyTrends(monthsBack int) ([]MonthlyTrend, error) {
	return a.trackingService.GetMonthlyTrends(monthsBack)
}

// GetVarianceReport returns planned vs actual analysis for projects in a date range
func (a *App) GetVarianceReport(startDate, endDate string) ([]VarianceReport, error) {
	return a.trackingService.GetVarianceReport(startDate, endDate)
}

// GetCapacityUtilization returns weekly capacity utilization data for a date range
func (a *App) GetCapacityUtilization(startDate, endDate string) ([]CapacityWeek, error) {
	return a.trackingService.GetCapacityUtilization(startDate, endDate)
}
