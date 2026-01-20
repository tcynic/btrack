package repository

import (
	"btrack/internal/database"
	"btrack/internal/models"
)

// NoteRepository handles note database operations
type NoteRepository struct {
	*Repository
}

// NewNoteRepository creates a new note repository
func NewNoteRepository(base *Repository) *NoteRepository {
	return &NoteRepository{Repository: base}
}

// GetByID retrieves a single note by ID
func (r *NoteRepository) GetByID(id int64) (*models.Note, error) {
	var note *models.Note
	err := r.QueryOne(database.SelectNoteByID, []any{id}, func(scan func(dest ...any) error) error {
		n, err := models.ScanNote(scan)
		if err != nil {
			return err
		}
		note = n
		return nil
	})
	if err != nil {
		return nil, err
	}
	return note, nil
}

// GetByProject retrieves all notes for a project
func (r *NoteRepository) GetByProject(projectID int64) ([]models.Note, error) {
	var notes []models.Note
	err := r.QuerySlice(database.SelectNotesByProject, []any{projectID}, func(scan func(dest ...any) error) error {
		n, err := models.ScanNote(scan)
		if err != nil {
			return err
		}
		notes = append(notes, *n)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if notes == nil {
		notes = []models.Note{}
	}
	return notes, nil
}

// Create inserts a new note
func (r *NoteRepository) Create(projectID int64, title, content string) (int64, error) {
	result, err := r.Exec(
		database.InsertNote,
		projectID, title, content,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update updates an existing note
func (r *NoteRepository) Update(id int64, title, content string) error {
	result, err := r.Exec(
		database.UpdateNote,
		title, content, id,
	)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("note")
	}
	
	return nil
}

// Delete removes a note
func (r *NoteRepository) Delete(id int64) error {
	result, err := r.Exec(database.DeleteNote, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("note")
	}
	
	return nil
}

// Search searches notes by title or content
func (r *NoteRepository) Search(query string) ([]models.NoteWithProject, error) {
	var notes []models.NoteWithProject
	err := r.QuerySlice(database.SearchNotes, []any{query, query}, func(scan func(dest ...any) error) error {
		n, err := models.ScanNoteWithProject(scan)
		if err != nil {
			return err
		}
		notes = append(notes, *n)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if notes == nil {
		notes = []models.NoteWithProject{}
	}
	return notes, nil
}
