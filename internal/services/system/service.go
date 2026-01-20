package system

import (
	"context"
	"database/sql"

	"btrack/internal/services/project"
)

// Service handles system-related operations (backup, export, templates)
type Service struct {
	projectService *project.Service
	db             *sql.DB
	ctx            context.Context
}

// NewService creates a new system service
func NewService(projectService *project.Service, db *sql.DB, ctx context.Context) *Service {
	return &Service{
		projectService: projectService,
		db:             db,
		ctx:            ctx,
	}
}
