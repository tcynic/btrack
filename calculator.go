package main

import (
	"time"

	"btrack/internal/models"
)

// CalculateDistribution calculates the frontloaded hour distribution
// without persisting to the database. Used for preview before saving.
func (a *App) CalculateDistribution(input models.CreateProjectInput) ([]models.WeeklyEntry, error) {
	return a.projectService.CalculateDistribution(input)
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
