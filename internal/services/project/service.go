package project

import (
	"context"
	"fmt"
	"time"

	"btrack/internal/models"
	"btrack/internal/store"
)

// Service handles project-related business logic
type Service struct {
	store *store.Store
}

// NewService creates a new project service
func NewService(s *store.Store) *Service {
	return &Service{
		store: s,
	}
}

// Create creates a new project with weekly entries
func (s *Service) Create(ctx context.Context, input models.CreateProjectInput) (*models.ProjectWithStats, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Calculate distribution
	entries, err := s.CalculateDistribution(input)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate distribution: %w", err)
	}

	// Create project
	project, err := s.store.CreateProject(input)
	if err != nil {
		return nil, err
	}

	// Add weekly entries
	for _, entry := range entries {
		_, err := s.store.AddWeeklyEntry(project.ID, entry)
		if err != nil {
			return nil, fmt.Errorf("failed to add weekly entry: %w", err)
		}
	}

	return s.GetByID(ctx, project.ID)
}

// GetByID retrieves a project with stats
func (s *Service) GetByID(ctx context.Context, id int64) (*models.ProjectWithStats, error) {
	p, err := s.store.GetProject(id)
	if err != nil {
		return nil, err
	}

	stats, err := s.getStats(p.ID)
	if err != nil {
		return nil, models.Internal(err, "failed to get project stats")
	}

	projectWithStats := &models.ProjectWithStats{
		Project:           p.Project,
		MyHours:           p.TotalSoldHours - p.SpecialistHours,
		TotalWeeks:        stats.TotalWeeks,
		TotalPlannedHours: stats.TotalPlanned,
		TotalActualHours:  stats.TotalActual,
	}
	projectWithStats.Health = s.CalculateHealth(*projectWithStats)
	return projectWithStats, nil
}

// GetAll retrieves all projects with stats
func (s *Service) GetAll(ctx context.Context, activeOnly bool) ([]models.ProjectWithStats, error) {
	projects, err := s.store.GetAllProjects(activeOnly)
	if err != nil {
		return nil, err
	}

	var result []models.ProjectWithStats
	for _, p := range projects {
		stats, err := s.getStats(p.ID)
		if err != nil {
			return nil, models.Internal(err, "failed to get project stats")
		}

		projectWithStats := models.ProjectWithStats{
			Project:           p.Project,
			MyHours:           p.TotalSoldHours - p.SpecialistHours,
			TotalWeeks:        stats.TotalWeeks,
			TotalPlannedHours: stats.TotalPlanned,
			TotalActualHours:  stats.TotalActual,
		}
		projectWithStats.Health = s.CalculateHealth(projectWithStats)
		result = append(result, projectWithStats)
	}

	if result == nil {
		result = []models.ProjectWithStats{}
	}
	return result, nil
}

// Search searches projects by name
func (s *Service) Search(ctx context.Context, query string) ([]models.Project, error) {
	return s.store.SearchProjects(query)
}

// getStats retrieves statistics for a project
func (s *Service) getStats(projectID int64) (struct{ TotalPlanned, TotalActual, TotalWeeks int }, error) {
	var stats struct {
		TotalPlanned int
		TotalActual  int
		TotalWeeks   int
	}

	totalPlanned, totalActual, totalWeeks, err := s.store.GetProjectStats(projectID)
	if err != nil {
		return stats, err
	}

	stats.TotalPlanned = totalPlanned
	stats.TotalActual = totalActual
	stats.TotalWeeks = totalWeeks
	return stats, nil
}

// CalculateHealth determines project health status
func (s *Service) CalculateHealth(p models.ProjectWithStats) models.ProjectHealth {
	// Persistent projects always on track
	if p.IsPersistent {
		return models.ProjectHealth{Status: "on_track"}
	}

	now := time.Now()
	endDate, _ := time.Parse("2006-01-02", p.EndDate)
	isCompleted := now.After(endDate)

	if isCompleted {
		if p.TotalActualHours <= p.TotalPlannedHours {
			return models.ProjectHealth{Status: "completed", Message: "Project completed successfully"}
		}
		return models.ProjectHealth{Status: "over_budget", Message: "Project completed over budget"}
	}

	if p.TotalActualHours == 0 {
		return models.ProjectHealth{Status: "on_track", Message: "No hours tracked yet"}
	}

	variance := float64(p.TotalActualHours) / float64(p.TotalPlannedHours)
	if variance > 1.2 {
		return models.ProjectHealth{Status: "over_budget", Message: "Significantly over planned hours"}
	} else if variance > 1.0 {
		return models.ProjectHealth{Status: "at_risk", Message: "Slightly over planned hours"}
	}

	return models.ProjectHealth{Status: "on_track"}
}

// CalculateDistribution calculates frontloaded hour distribution
func (s *Service) CalculateDistribution(input models.CreateProjectInput) ([]models.WeeklyEntry, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	myHours := input.TotalSoldHours - input.SpecialistHours

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, models.ValidationError("startDate", "invalid date format")
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, models.ValidationError("endDate", "invalid date format")
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
			hours++ // Frontload remainder to earliest weeks
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

// Helper functions

func calculateWeeksBetween(start, end time.Time) int {
	days := int(end.Sub(start).Hours()/24) + 1
	weeks := days / 7
	if days%7 > 0 {
		weeks++
	}
	if weeks < 1 {
		weeks = 1
	}
	return weeks
}

func getMonday(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return t.AddDate(0, 0, -(weekday - 1))
}
