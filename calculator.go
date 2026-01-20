package main

import (
	"btrack/internal/models"
)

// CalculateDistribution calculates the frontloaded hour distribution
// without persisting to the database. Used for preview before saving.
func (a *App) CalculateDistribution(input models.CreateProjectInput) ([]models.WeeklyEntry, error) {
	return a.projectService.CalculateDistribution(input)
}
