package notes

import (
	"fmt"

	"btrack/internal/database"
	"btrack/internal/models"
)

// CreateMeeting creates a new meeting for a project
func (s *Service) CreateMeeting(input models.CreateMeetingInput) (*models.Meeting, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	meetingID, err := s.meetingRepo.Create(
		input.ProjectID,
		input.Title,
		input.MeetingDate,
		input.DurationMinutes,
		input.Attendees,
		input.Notes,
	)
	if err != nil {
		return nil, err
	}

	return s.GetMeeting(meetingID)
}

// GetMeetings returns all meetings for a project
func (s *Service) GetMeetings(projectID int64) ([]models.Meeting, error) {
	meetings, err := s.meetingRepo.GetByProject(projectID)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(meetings), nil
}

// GetMeeting returns a single meeting by ID
func (s *Service) GetMeeting(id int64) (*models.Meeting, error) {
	return s.meetingRepo.GetByID(id)
}

// UpdateMeeting updates an existing meeting
func (s *Service) UpdateMeeting(input models.UpdateMeetingInput) (*models.Meeting, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	err := s.meetingRepo.Update(
		input.ID,
		input.Title,
		input.MeetingDate,
		input.DurationMinutes,
		input.Attendees,
		input.Notes,
	)
	if err != nil {
		return nil, err
	}

	return s.GetMeeting(input.ID)
}

// DeleteMeeting removes a meeting (also deletes associated tasks)
func (s *Service) DeleteMeeting(id int64) error {
	// First delete associated tasks
	if err := s.taskRepo.DeleteBySource("meeting", id); err != nil {
		return err
	}

	return s.meetingRepo.Delete(id)
}

// GetMeetingsByDate returns all meetings for a specific date across all projects
func (s *Service) GetMeetingsByDate(date string) ([]models.MeetingWithProject, error) {
	rows, err := s.db.Query(database.SelectMeetingsByDate, date)
	if err != nil {
		return nil, fmt.Errorf("failed to query meetings by date: %w", err)
	}
	defer rows.Close()

	var meetings []models.MeetingWithProject
	for rows.Next() {
		m, err := models.ScanMeetingWithProject(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan meeting: %w", err)
		}
		meetings = append(meetings, *m)
	}

	return database.EnsureSlice(meetings), nil
}

// GetMeetingsByWeek returns all meetings for a week (7-day period starting from weekStartDate)
func (s *Service) GetMeetingsByWeek(weekStartDate string) ([]models.MeetingWithProject, error) {
	// Calculate the end date (7 days after start)
	startDate, err := database.ParseDate(weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid week start date: %w", err)
	}
	endDate := startDate.AddDate(0, 0, 7)

	rows, err := s.db.Query(database.SelectMeetingsByWeek, weekStartDate, endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to query meetings by week: %w", err)
	}
	defer rows.Close()

	var meetings []models.MeetingWithProject
	for rows.Next() {
		m, err := models.ScanMeetingWithProject(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("failed to scan meeting: %w", err)
		}
		meetings = append(meetings, *m)
	}

	return database.EnsureSlice(meetings), nil
}

// SearchMeetings searches for meetings by title, notes, or attendees
func (s *Service) SearchMeetings(query string) ([]models.MeetingWithProject, error) {
	searchPattern := "%" + query + "%"
	meetings, err := s.meetingRepo.Search(searchPattern)
	if err != nil {
		return nil, err
	}
	return database.EnsureSlice(meetings), nil
}
