package main

import (
	"context"
	"database/sql"
	"log"

	"btrack/internal/database"
	"btrack/internal/repository"
	"btrack/internal/services/notes"
	"btrack/internal/services/project"
	"btrack/internal/services/system"
	"btrack/internal/services/tracking"
)

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
	
	// Repositories
	projects      *repository.ProjectRepository
	weeklyEntries *repository.WeeklyEntryRepository
	meetings      *repository.MeetingRepository
	notes         *repository.NoteRepository
	goals         *repository.GoalRepository
	tasks         *repository.TaskRepository
	
	// Services
	projectService  *project.Service
	trackingService *tracking.Service
	notesService    *notes.Service
	systemService   *system.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	db, err := database.Initialize()
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return
	}
	a.db = db
	
	// Initialize repositories
	base := repository.NewRepository(db)
	a.projects = repository.NewProjectRepository(base)
	a.weeklyEntries = repository.NewWeeklyEntryRepository(base)
	a.meetings = repository.NewMeetingRepository(base)
	a.notes = repository.NewNoteRepository(base)
	a.goals = repository.NewGoalRepository(base)
	a.tasks = repository.NewTaskRepository(base)
	
	// Initialize services
	a.projectService = project.NewService(a.projects, a.weeklyEntries, db)
	a.trackingService = tracking.NewService(a.weeklyEntries, a.projects, a.goals, db)
	a.notesService = notes.NewService(a.meetings, a.notes, a.goals, a.tasks, db)
	a.systemService = system.NewService(a.projectService, db, ctx)
	
	log.Println("Database initialized successfully")
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
		log.Println("Database connection closed")
	}
}
