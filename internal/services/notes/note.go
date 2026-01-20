package notes

import (
	"btrack/internal/database"
	"btrack/internal/models"
)

// CreateNote creates a new note for a project
func (s *Service) CreateNote(input models.CreateNoteInput) (*models.Note, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	noteID, err := s.noteRepo.Create(input.ProjectID, input.Title, input.Content)
	if err != nil {
		return nil, err
	}

	return s.GetNote(noteID)
}

// GetNotes returns all notes for a project
func (s *Service) GetNotes(projectID int64) ([]models.Note, error) {
	notes, err := s.noteRepo.GetByProject(projectID)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(notes), nil
}

// GetNote returns a single note by ID
func (s *Service) GetNote(id int64) (*models.Note, error) {
	return s.noteRepo.GetByID(id)
}

// UpdateNote updates an existing note
func (s *Service) UpdateNote(input models.UpdateNoteInput) (*models.Note, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	err := s.noteRepo.Update(input.ID, input.Title, input.Content)
	if err != nil {
		return nil, err
	}

	return s.GetNote(input.ID)
}

// DeleteNote removes a note (also deletes associated tasks)
func (s *Service) DeleteNote(id int64) error {
	// First delete associated tasks
	if err := s.taskRepo.DeleteBySource("note", id); err != nil {
		return err
	}

	return s.noteRepo.Delete(id)
}

// SearchNotes searches for notes by title or content
func (s *Service) SearchNotes(query string) ([]models.NoteWithProject, error) {
	searchPattern := "%" + query + "%"
	notes, err := s.noteRepo.Search(searchPattern)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(notes), nil
}
