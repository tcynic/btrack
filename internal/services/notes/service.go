package notes

import (
	"btrack/internal/store"
)

// Service handles notes-related business logic (meetings, notes, goals, tasks)
type Service struct {
	store *store.Store
}

// NewService creates a new notes service
func NewService(s *store.Store) *Service {
	return &Service{
		store: s,
	}
}
