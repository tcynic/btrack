package notes

import (
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

	goal := &models.Goal{
		ProjectID:   input.ProjectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      models.GoalStatusPending,
		TargetDate:  input.TargetDate,
	}

	if err := s.store.AddGoal(input.ProjectID, goal); err != nil {
		return nil, err
	}

	return goal, nil
}

// GetGoals returns all goals for a project
func (s *Service) GetGoals(projectID int64) ([]models.Goal, error) {
	return s.store.GetGoals(projectID)
}

// GetGoal returns a single goal by ID
func (s *Service) GetGoal(id int64) (*models.Goal, error) {
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, g := range proj.Goals {
			if g.ID == id {
				return &g, nil
			}
		}
	}
	return nil, models.NotFound("goal")
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

	// Find the project this goal belongs to
	var projectID int64
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, g := range proj.Goals {
			if g.ID == input.ID {
				projectID = proj.ID
				break
			}
		}
		if projectID > 0 {
			break
		}
	}

	if projectID == 0 {
		return nil, models.NotFound("goal")
	}

	goal := &models.Goal{
		ID:          input.ID,
		ProjectID:   projectID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		TargetDate:  input.TargetDate,
	}

	if err := s.store.UpdateGoal(projectID, goal); err != nil {
		return nil, err
	}

	return goal, nil
}

// UpdateGoalStatus updates only the status of a goal
func (s *Service) UpdateGoalStatus(id int64, status string) (*models.Goal, error) {
	if !models.IsValidStatus(status) {
		return nil, models.ValidationError("status", "invalid status value")
	}

	// Get the goal first
	goal, err := s.GetGoal(id)
	if err != nil {
		return nil, err
	}

	// Update status
	goal.Status = status

	// Find project and update
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, g := range proj.Goals {
			if g.ID == id {
				if err := s.store.UpdateGoal(proj.ID, goal); err != nil {
					return nil, err
				}
				return goal, nil
			}
		}
	}

	return nil, models.NotFound("goal")
}

// DeleteGoal removes a goal
func (s *Service) DeleteGoal(id int64) error {
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, g := range proj.Goals {
			if g.ID == id {
				return s.store.DeleteGoal(proj.ID, id)
			}
		}
	}
	return models.NotFound("goal")
}

// GetGoalStats returns statistics for a project's goals
func (s *Service) GetGoalStats(projectID int64) (*GoalStats, error) {
	goals, err := s.store.GetGoals(projectID)
	if err != nil {
		return nil, err
	}

	stats := &GoalStats{}
	for _, goal := range goals {
		stats.Total++
		switch goal.Status {
		case models.GoalStatusPending:
			stats.Pending++
		case models.GoalStatusInProgress:
			stats.InProgress++
		case models.GoalStatusCompleted:
			stats.Completed++
		case models.GoalStatusCancelled:
			stats.Cancelled++
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
	return s.store.SearchAllGoals(query)
}
