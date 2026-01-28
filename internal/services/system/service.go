package system

import (
	"context"

	"btrack/internal/services/project"
	"btrack/internal/store"
)

// Service handles system-related operations (backup, export, templates)
type Service struct {
	projectService *project.Service
	store          *store.Store
	ctx            context.Context
}

// NewService creates a new system service
func NewService(projectService *project.Service, s *store.Store, ctx context.Context) *Service {
	return &Service{
		projectService: projectService,
		store:          s,
		ctx:            ctx,
	}
}
