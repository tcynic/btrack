package models

// Note represents a project note with markdown content
type Note struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"projectId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// NoteWithProject includes the parent project name for search results
type NoteWithProject struct {
	Note
	ProjectName string `json:"projectName"`
}

// ScanNote scans a database row into a Note struct.
// Expected columns: id, project_id, title, content, created_at, updated_at
func ScanNote(scan func(dest ...any) error) (*Note, error) {
	var n Note
	var content *string
	var createdAt, updatedAt string

	err := scan(
		&n.ID,
		&n.ProjectID,
		&n.Title,
		&content,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	n.Content = NullableString(content)
	n.CreatedAt = createdAt
	n.UpdatedAt = updatedAt

	return &n, nil
}

// ScanNoteWithProject scans a database row into a NoteWithProject struct.
// Expected columns: id, project_id, title, content, created_at, updated_at, project_name
func ScanNoteWithProject(scan func(dest ...any) error) (*NoteWithProject, error) {
	var n NoteWithProject
	var content *string
	var createdAt, updatedAt string

	err := scan(
		&n.ID,
		&n.ProjectID,
		&n.Title,
		&content,
		&createdAt,
		&updatedAt,
		&n.ProjectName,
	)
	if err != nil {
		return nil, err
	}

	n.Content = NullableString(content)
	n.CreatedAt = createdAt
	n.UpdatedAt = updatedAt

	return &n, nil
}

// CreateNoteInput is the input for creating a new note
type CreateNoteInput struct {
	ProjectID int64  `json:"projectId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}

// UpdateNoteInput is the input for updating a note
type UpdateNoteInput struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Validate checks if the CreateNoteInput is valid
func (c *CreateNoteInput) Validate() error {
	if c.Title == "" {
		return ValidationError("title", "title is required")
	}
	return nil
}

// Validate checks if the UpdateNoteInput is valid
func (u *UpdateNoteInput) Validate() error {
	if u.Title == "" {
		return ValidationError("title", "title is required")
	}
	return nil
}
