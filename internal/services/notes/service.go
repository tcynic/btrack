package notes

import (
	"database/sql"

	"btrack/internal/repository"
)

// Service handles notes-related business logic (meetings, notes, goals, tasks)
type Service struct {
	meetingRepo *repository.MeetingRepository
	noteRepo    *repository.NoteRepository
	goalRepo    *repository.GoalRepository
	taskRepo    *repository.TaskRepository
	db          *sql.DB
}

// NewService creates a new notes service
func NewService(
	meetingRepo *repository.MeetingRepository,
	noteRepo *repository.NoteRepository,
	goalRepo *repository.GoalRepository,
	taskRepo *repository.TaskRepository,
	db *sql.DB,
) *Service {
	return &Service{
		meetingRepo: meetingRepo,
		noteRepo:    noteRepo,
		goalRepo:    goalRepo,
		taskRepo:    taskRepo,
		db:          db,
	}
}
