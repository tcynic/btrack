package main

import (
	"fmt"
	"time"

	"btrack/internal/models"
)

// CalculateDistribution calculates the frontloaded hour distribution
// without persisting to the database. Used for preview before saving.
//
// Algorithm:
// 1. MyHours = TotalSoldHours - SpecialistHours
// 2. Calculate number of weeks between StartDate and EndDate
// 3. BaseHours = MyHours / Weeks (integer division)
// 4. Remainder = MyHours % Weeks
// 5. Assign BaseHours to every week
// 6. Add 1 hour to earliest weeks until remainder exhausted
//
// Example: 10 hours / 3 weeks → Week1: 4, Week2: 3, Week3: 3
func (a *App) CalculateDistribution(input models.CreateProjectInput) ([]models.WeeklyEntry, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	myHours := input.TotalSoldHours - input.SpecialistHours

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	weeks := calculateWeeksBetween(startDate, endDate)
	if weeks < 1 {
		weeks = 1
	}

	baseHours := myHours / weeks
	remainder := myHours % weeks

	entries := make([]models.WeeklyEntry, weeks)
	currentMonday := getMonday(startDate)

	for i := 0; i < weeks; i++ {
		hours := baseHours
		if i < remainder {
			hours++ // Add 1 to earliest weeks (frontloading)
		}

		entries[i] = models.WeeklyEntry{
			WeekStartDate: currentMonday.Format("2006-01-02"),
			WeekNumber:    i + 1,
			PlannedHours:  hours,
			ActualHours:   0,
		}

		currentMonday = currentMonday.AddDate(0, 0, 7)
	}

	return entries, nil
}

// calculateWeeksBetween calculates the number of weeks between two dates
func calculateWeeksBetween(start, end time.Time) int {
	days := int(end.Sub(start).Hours()/24) + 1
	weeks := days / 7
	if days%7 > 0 {
		weeks++ // Round up partial weeks
	}
	if weeks < 1 {
		weeks = 1
	}
	return weeks
}

// getMonday returns the Monday of the week containing the given date
func getMonday(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	return t.AddDate(0, 0, -(weekday - 1))
}

// getCurrentWeekMonday returns the Monday of the current week
func getCurrentWeekMonday() time.Time {
	return getMonday(time.Now())
}
