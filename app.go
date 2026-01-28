package main

import (
	"context"
	"log"

	"btrack/internal/services/notes"
	"btrack/internal/services/project"
	"btrack/internal/services/system"
	"btrack/internal/services/tracking"
	"btrack/internal/store"
)

// App struct
type App struct {
	ctx   context.Context
	store *store.Store
	
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

	// Initialize store
	a.store = store.NewStore()
	if err := a.store.Load(); err != nil {
		log.Printf("Failed to load store: %v", err)
		return
	}
	
	// Initialize services with store
	a.projectService = project.NewService(a.store)
	a.trackingService = tracking.NewService(a.store)
	a.notesService = notes.NewService(a.store)
	a.systemService = system.NewService(a.projectService, a.store, ctx)
	
	log.Println("Store initialized successfully")
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.store != nil {
		if err := a.store.Save(); err != nil {
			log.Printf("Failed to save store on shutdown: %v", err)
		} else {
			log.Println("Store saved successfully")
		}
	}
}
