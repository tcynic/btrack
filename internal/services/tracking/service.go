package tracking

import (
	"btrack/internal/store"
)

// Service handles tracking-related business logic
type Service struct {
	store *store.Store
}

// NewService creates a new tracking service
func NewService(s *store.Store) *Service {
	return &Service{
		store: s,
	}
}
