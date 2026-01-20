package main

import (
	"btrack/internal/database"
	"btrack/internal/models"
)

// CreateNote creates a new note for a project
func (a *App) CreateNote(input models.CreateNoteInput) (*models.Note, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	noteID, err := a.notes.Create(input.ProjectID, input.Title, input.Content)
	if err != nil {
		return nil, err
	}

	return a.GetNote(noteID)
}

// GetNotes returns all notes for a project
func (a *App) GetNotes(projectID int64) ([]models.Note, error) {
	notes, err := a.notes.GetByProject(projectID)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(notes), nil
}

// GetNote returns a single note by ID
func (a *App) GetNote(id int64) (*models.Note, error) {
	return a.notes.GetByID(id)
}

// UpdateNote updates an existing note
func (a *App) UpdateNote(input models.UpdateNoteInput) (*models.Note, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	err := a.notes.Update(input.ID, input.Title, input.Content)
	if err != nil {
		return nil, err
	}

	return a.GetNote(input.ID)
}

// DeleteNote removes a note
func (a *App) DeleteNote(id int64) error {
	// First delete associated tasks
	if err := a.tasks.DeleteBySource("note", id); err != nil {
		return err
	}

	return a.notes.Delete(id)
}

// SearchNotes searches for notes by title or content
func (a *App) SearchNotes(query string) ([]models.NoteWithProject, error) {
	searchPattern := "%" + query + "%"
	notes, err := a.notes.Search(searchPattern)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(notes), nil
}
