package main

import (
	"fmt"

	"btrack/internal/database"
	"btrack/internal/models"
)

// CreateNote creates a new note for a project
func (a *App) CreateNote(input models.CreateNoteInput) (*models.Note, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	result, err := a.db.Exec(database.InsertNote,
		input.ProjectID,
		input.Title,
		input.Content,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert note: %w", err)
	}

	noteID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get note ID: %w", err)
	}

	return a.GetNote(noteID)
}

// GetNotes returns all notes for a project
func (a *App) GetNotes(projectID int64) ([]models.Note, error) {
	rows, err := a.db.Query(database.SelectNotesByProject, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		n, err := models.ScanNote(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, *n)
	}

	return database.EnsureSlice(notes), nil
}

// GetNote returns a single note by ID
func (a *App) GetNote(id int64) (*models.Note, error) {
	row := a.db.QueryRow(database.SelectNoteByID, id)
	n, err := models.ScanNote(row.Scan)
	if err != nil {
		return nil, models.ErrNoteNotFound
	}
	return n, nil
}

// UpdateNote updates an existing note
func (a *App) UpdateNote(input models.UpdateNoteInput) (*models.Note, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	result, err := a.db.Exec(database.UpdateNote,
		input.Title,
		input.Content,
		input.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, models.ErrNoteNotFound
	}

	return a.GetNote(input.ID)
}

// DeleteNote removes a note
func (a *App) DeleteNote(id int64) error {
	// First delete associated tasks
	_, err := a.db.Exec(database.DeleteTasksBySource, "note", id)
	if err != nil {
		return fmt.Errorf("failed to delete associated tasks: %w", err)
	}

	result, err := a.db.Exec(database.DeleteNote, id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return models.ErrNoteNotFound
	}

	return nil
}

// SearchNotes searches for notes by title or content
func (a *App) SearchNotes(query string) ([]models.NoteWithProject, error) {
	rows, err := a.db.Query(database.SearchNotes, query, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search notes: %w", err)
	}
	defer rows.Close()

	var notes []models.NoteWithProject
	for rows.Next() {
		n, err := models.ScanNoteWithProject(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, *n)
	}

	return database.EnsureSlice(notes), nil
}
