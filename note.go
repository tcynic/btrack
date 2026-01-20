package main

import (
	"btrack/internal/models"
)

// CreateNote creates a new note for a project
func (a *App) CreateNote(input models.CreateNoteInput) (*models.Note, error) {
	return a.notesService.CreateNote(input)
}

// GetNotes returns all notes for a project
func (a *App) GetNotes(projectID int64) ([]models.Note, error) {
	return a.notesService.GetNotes(projectID)
}

// GetNote returns a single note by ID
func (a *App) GetNote(id int64) (*models.Note, error) {
	return a.notesService.GetNote(id)
}

// UpdateNote updates an existing note
func (a *App) UpdateNote(input models.UpdateNoteInput) (*models.Note, error) {
	return a.notesService.UpdateNote(input)
}

// DeleteNote removes a note
func (a *App) DeleteNote(id int64) error {
	return a.notesService.DeleteNote(id)
}

// SearchNotes searches for notes by title or content
func (a *App) SearchNotes(query string) ([]models.NoteWithProject, error) {
	return a.notesService.SearchNotes(query)
}
