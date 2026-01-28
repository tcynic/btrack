package notes

import (
	"btrack/internal/models"
)

// CreateNote creates a new note for a project
func (s *Service) CreateNote(input models.CreateNoteInput) (*models.Note, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	note := &models.Note{
		ProjectID: input.ProjectID,
		Title:     input.Title,
		Content:   input.Content,
	}

	if err := s.store.AddNote(input.ProjectID, note); err != nil {
		return nil, err
	}

	return note, nil
}

// GetNotes returns all notes for a project
func (s *Service) GetNotes(projectID int64) ([]models.Note, error) {
	return s.store.GetNotes(projectID)
}

// GetNote returns a single note by ID
func (s *Service) GetNote(id int64) (*models.Note, error) {
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, n := range proj.Notes {
			if n.ID == id {
				return &n, nil
			}
		}
	}
	return nil, models.NotFound("note")
}

// UpdateNote updates an existing note
func (s *Service) UpdateNote(input models.UpdateNoteInput) (*models.Note, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Find the project this note belongs to
	var projectID int64
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, n := range proj.Notes {
			if n.ID == input.ID {
				projectID = proj.ID
				break
			}
		}
		if projectID > 0 {
			break
		}
	}

	if projectID == 0 {
		return nil, models.NotFound("note")
	}

	note := &models.Note{
		ID:        input.ID,
		ProjectID: projectID,
		Title:     input.Title,
		Content:   input.Content,
	}

	if err := s.store.UpdateNote(projectID, note); err != nil {
		return nil, err
	}

	return note, nil
}

// DeleteNote removes a note (also deletes associated tasks)
func (s *Service) DeleteNote(id int64) error {
	// Find which project this note belongs to
	var projectID int64
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, n := range proj.Notes {
			if n.ID == id {
				projectID = proj.ID
				// Delete associated tasks first
				for _, task := range proj.Tasks {
					if task.SourceType == "note" && task.SourceID == id {
						s.store.DeleteTask(proj.ID, task.ID)
					}
				}
				break
			}
		}
		if projectID > 0 {
			break
		}
	}

	if projectID == 0 {
		return models.NotFound("note")
	}

	return s.store.DeleteNote(projectID, id)
}

// SearchNotes searches for notes by title or content
func (s *Service) SearchNotes(query string) ([]models.NoteWithProject, error) {
	return s.store.SearchAllNotes(query)
}
