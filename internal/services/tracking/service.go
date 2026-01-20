package tracking

import (
	"database/sql"

	"btrack/internal/repository"
)

// Service handles tracking-related business logic
type Service struct {
	weeklyRepo  *repository.WeeklyEntryRepository
	projectRepo *repository.ProjectRepository
	goalRepo    *repository.GoalRepository
	db          *sql.DB
}

// NewService creates a new tracking service
func NewService(
	weeklyRepo *repository.WeeklyEntryRepository,
	projectRepo *repository.ProjectRepository,
	goalRepo *repository.GoalRepository,
	db *sql.DB,
) *Service {
	return &Service{
		weeklyRepo:  weeklyRepo,
		projectRepo: projectRepo,
		goalRepo:    goalRepo,
		db:          db,
	}
}
