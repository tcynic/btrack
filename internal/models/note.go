package models

import "time"

// Note represents a project note with markdown content
type Note struct {
	ID        int64     `json:"id"`
	ProjectID int64     `json:"projectId"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
		return ErrTitleRequired
	}
	return nil
}

// Validate checks if the UpdateNoteInput is valid
func (u *UpdateNoteInput) Validate() error {
	if u.Title == "" {
		return ErrTitleRequired
	}
	return nil
}
