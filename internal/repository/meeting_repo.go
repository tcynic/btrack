package repository

import (
	"btrack/internal/database"
	"btrack/internal/models"
)

// MeetingRepository handles meeting database operations
type MeetingRepository struct {
	*Repository
}

// NewMeetingRepository creates a new meeting repository
func NewMeetingRepository(base *Repository) *MeetingRepository {
	return &MeetingRepository{Repository: base}
}

// GetByID retrieves a single meeting by ID
func (r *MeetingRepository) GetByID(id int64) (*models.Meeting, error) {
	var meeting *models.Meeting
	err := r.QueryOne(database.SelectMeetingByID, []any{id}, func(scan func(dest ...any) error) error {
		m, err := models.ScanMeeting(scan)
		if err != nil {
			return err
		}
		meeting = m
		return nil
	})
	if err != nil {
		return nil, err
	}
	return meeting, nil
}

// GetByProject retrieves all meetings for a project
func (r *MeetingRepository) GetByProject(projectID int64) ([]models.Meeting, error) {
	var meetings []models.Meeting
	err := r.QuerySlice(database.SelectMeetingsByProject, []any{projectID}, func(scan func(dest ...any) error) error {
		m, err := models.ScanMeeting(scan)
		if err != nil {
			return err
		}
		meetings = append(meetings, *m)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if meetings == nil {
		meetings = []models.Meeting{}
	}
	return meetings, nil
}

// Create inserts a new meeting
func (r *MeetingRepository) Create(projectID int64, title, meetingDate string, durationMinutes int, attendees, notes string) (int64, error) {
	result, err := r.Exec(
		database.InsertMeeting,
		projectID, title, meetingDate, durationMinutes, attendees, notes,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update updates an existing meeting
func (r *MeetingRepository) Update(id int64, title, meetingDate string, durationMinutes int, attendees, notes string) error {
	result, err := r.Exec(
		database.UpdateMeeting,
		title, meetingDate, durationMinutes, attendees, notes, id,
	)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("meeting")
	}
	
	return nil
}

// Delete removes a meeting
func (r *MeetingRepository) Delete(id int64) error {
	result, err := r.Exec(database.DeleteMeeting, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.DatabaseError(err)
	}
	
	if rowsAffected == 0 {
		return models.NotFound("meeting")
	}
	
	return nil
}

// Search searches meetings by title, notes, or attendees
func (r *MeetingRepository) Search(query string) ([]models.MeetingWithProject, error) {
	var meetings []models.MeetingWithProject
	err := r.QuerySlice(database.SearchMeetings, []any{query, query, query}, func(scan func(dest ...any) error) error {
		m, err := models.ScanMeetingWithProject(scan)
		if err != nil {
			return err
		}
		meetings = append(meetings, *m)
		return nil
	})
	if err != nil {
		return nil, err
	}

	if meetings == nil {
		meetings = []models.MeetingWithProject{}
	}
	return meetings, nil
}
