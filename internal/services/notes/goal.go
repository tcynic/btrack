package notes

import (
	"fmt"

	"btrack/internal/database"
	"btrack/internal/models"
)

// GoalStats represents statistics for a project's goals
type GoalStats struct {
	Total          int     `json:"total"`
	Pending        int     `json:"pending"`
	InProgress     int     `json:"inProgress"`
	Completed      int     `json:"completed"`
	Cancelled      int     `json:"cancelled"`
	CompletionRate float64 `json:"completionRate"` // percentage
}

// CreateGoal creates a new goal for a project
func (s *Service) CreateGoal(input models.CreateGoalInput) (*models.Goal, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	goalID, err := s.goalRepo.Create(input.ProjectID, input.Title, input.Description, input.TargetDate)
	if err != nil {
		return nil, err
	}

	return s.GetGoal(goalID)
}

// GetGoals returns all goals for a project
func (s *Service) GetGoals(projectID int64) ([]models.Goal, error) {
	goals, err := s.goalRepo.GetByProject(projectID)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(goals), nil
}

// GetGoal returns a single goal by ID
func (s *Service) GetGoal(id int64) (*models.Goal, error) {
	return s.goalRepo.GetByID(id)
}

// UpdateGoal updates an existing goal
func (s *Service) UpdateGoal(input models.UpdateGoalInput) (*models.Goal, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Default status if not provided
	if input.Status == "" {
		input.Status = models.GoalStatusPending
	}

	err := s.goalRepo.Update(input.ID, input.Title, input.Description, input.Status, input.TargetDate)
	if err != nil {
		return nil, err
	}

	return s.GetGoal(input.ID)
}

// UpdateGoalStatus updates only the status of a goal
func (s *Service) UpdateGoalStatus(id int64, status string) (*models.Goal, error) {
	if !models.IsValidStatus(status) {
		return nil, models.ValidationError("status", "invalid status value")
	}

	err := s.goalRepo.UpdateStatus(id, status)
	if err != nil {
		return nil, err
	}

	return s.GetGoal(id)
}

// DeleteGoal removes a goal
func (s *Service) DeleteGoal(id int64) error {
	return s.goalRepo.Delete(id)
}

// GetGoalStats returns statistics for a project's goals
func (s *Service) GetGoalStats(projectID int64) (*GoalStats, error) {
	rows, err := s.db.Query(database.SelectGoalStatsByProject, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query goal stats: %w", err)
	}
	defer rows.Close()

	stats := &GoalStats{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan goal stat: %w", err)
		}

		stats.Total += count
		switch status {
		case models.GoalStatusPending:
			stats.Pending = count
		case models.GoalStatusInProgress:
			stats.InProgress = count
		case models.GoalStatusCompleted:
			stats.Completed = count
		case models.GoalStatusCancelled:
			stats.Cancelled = count
		}
	}

	// Calculate completion rate
	if stats.Total > 0 {
		stats.CompletionRate = float64(stats.Completed) / float64(stats.Total) * 100
	}

	return stats, nil
}

// SearchGoals searches for goals by title or description
func (s *Service) SearchGoals(query string) ([]models.GoalWithProject, error) {
	searchPattern := "%" + query + "%"
	goals, err := s.goalRepo.Search(searchPattern)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(goals), nil
}
