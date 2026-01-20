package repository

import (
	"btrack/internal/database"
	"btrack/internal/models"
)

// GoalRepository handles goal database operations
type GoalRepository struct {
	*Repository
}

// NewGoalRepository creates a new goal repository
func NewGoalRepository(base *Repository) *GoalRepository {
	return &GoalRepository{Repository: base}
}

// GetByID retrieves a single goal by ID
func (r *GoalRepository) GetByID(id int64) (*models.Goal, error) {
	var goal *models.Goal
	err := r.QueryOne(database.SelectGoalByID, []any{id}, func(scan func(dest ...any) error) error {
		g, err := models.ScanGoal(scan)
		if err != nil {
			return err
		}
		goal = g
		return nil
	})
	if err != nil {
		return nil, err
	}
	return goal, nil
}

// GetByProject retrieves all goals for a project
func (r *GoalRepository) GetByProject(projectID int64) ([]models.Goal, error) {
	var goals []models.Goal
	err := r.QuerySlice(database.SelectGoalsByProject, []any{projectID}, func(scan func(dest ...any) error) error {
		g, err := models.ScanGoal(scan)
		if err != nil {
			return err
		}
		goals = append(goals, *g)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if goals == nil {
		goals = []models.Goal{}
	}
	return goals, nil
}

// Create inserts a new goal
func (r *GoalRepository) Create(projectID int64, title, description, targetDate string) (int64, error) {
	result, err := r.Exec(
		database.InsertGoal,
		projectID, title, description, targetDate,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update updates an existing goal
func (r *GoalRepository) Update(id int64, title, description, status, targetDate string) error {
	result, err := r.Exec(
		database.UpdateGoal,
		title, description, status, targetDate, id,
	)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("goal")
	}
	
	return nil
}

// UpdateStatus updates only the status of a goal
func (r *GoalRepository) UpdateStatus(id int64, status string) error {
	result, err := r.Exec(database.UpdateGoalStatus, status, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("goal")
	}
	
	return nil
}

// Delete removes a goal
func (r *GoalRepository) Delete(id int64) error {
	result, err := r.Exec(database.DeleteGoal, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("goal")
	}
	
	return nil
}

// Search searches goals by title or description
func (r *GoalRepository) Search(query string) ([]models.GoalWithProject, error) {
	var goals []models.GoalWithProject
	err := r.QuerySlice(database.SearchGoals, []any{query, query}, func(scan func(dest ...any) error) error {
		g, err := models.ScanGoalWithProject(scan)
		if err != nil {
			return err
		}
		goals = append(goals, *g)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if goals == nil {
		goals = []models.GoalWithProject{}
	}
	return goals, nil
}
