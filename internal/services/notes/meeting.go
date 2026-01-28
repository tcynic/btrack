package notes

import (
	"fmt"
	"time"

	"btrack/internal/models"
)

// CreateMeeting creates a new meeting for a project
func (s *Service) CreateMeeting(input models.CreateMeetingInput) (*models.Meeting, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	meeting := &models.Meeting{
		ProjectID:       input.ProjectID,
		Title:           input.Title,
		MeetingDate:     input.MeetingDate,
		DurationMinutes: input.DurationMinutes,
		Attendees:       input.Attendees,
		Notes:           input.Notes,
	}

	if err := s.store.AddMeeting(input.ProjectID, meeting); err != nil {
		return nil, err
	}

	return meeting, nil
}

// GetMeetings returns all meetings for a project
func (s *Service) GetMeetings(projectID int64) ([]models.Meeting, error) {
	return s.store.GetMeetings(projectID)
}

// GetMeeting returns a single meeting by ID
func (s *Service) GetMeeting(id int64) (*models.Meeting, error) {
	// Need to find which project this meeting belongs to
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, m := range proj.Meetings {
			if m.ID == id {
				return &m, nil
			}
		}
	}
	return nil, models.NotFound("meeting")
}

// UpdateMeeting updates an existing meeting
func (s *Service) UpdateMeeting(input models.UpdateMeetingInput) (*models.Meeting, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Find the project this meeting belongs to
	var projectID int64
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, m := range proj.Meetings {
			if m.ID == input.ID {
				projectID = proj.ID
				break
			}
		}
		if projectID > 0 {
			break
		}
	}

	if projectID == 0 {
		return nil, models.NotFound("meeting")
	}

	meeting := &models.Meeting{
		ID:              input.ID,
		ProjectID:       projectID,
		Title:           input.Title,
		MeetingDate:     input.MeetingDate,
		DurationMinutes: input.DurationMinutes,
		Attendees:       input.Attendees,
		Notes:           input.Notes,
	}

	if err := s.store.UpdateMeeting(projectID, meeting); err != nil {
		return nil, err
	}

	return meeting, nil
}

// DeleteMeeting removes a meeting (also deletes associated tasks)
func (s *Service) DeleteMeeting(id int64) error {
	// Find which project this meeting belongs to
	var projectID int64
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, m := range proj.Meetings {
			if m.ID == id {
				projectID = proj.ID
				// Delete associated tasks first
				for _, task := range proj.Tasks {
					if task.SourceType == "meeting" && task.SourceID == id {
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
		return models.NotFound("meeting")
	}

	return s.store.DeleteMeeting(projectID, id)
}

// GetMeetingsByDate returns all meetings for a specific date across all projects
func (s *Service) GetMeetingsByDate(date string) ([]models.MeetingWithProject, error) {
	projects, err := s.store.GetAllProjects(false)
	if err != nil {
		return nil, err
	}

	var meetings []models.MeetingWithProject
	for _, proj := range projects {
		for _, m := range proj.Meetings {
			if m.MeetingDate == date {
				meetings = append(meetings, models.MeetingWithProject{
					Meeting:     m,
					ProjectName: proj.Name,
				})
			}
		}
	}

	return meetings, nil
}

// GetMeetingsByWeek returns all meetings for a week (7-day period starting from weekStartDate)
func (s *Service) GetMeetingsByWeek(weekStartDate string) ([]models.MeetingWithProject, error) {
	// Calculate the end date (7 days after start)
	startDate, err := time.Parse("2006-01-02", weekStartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid week start date: %w", err)
	}
	endDate := startDate.AddDate(0, 0, 7)

	projects, err := s.store.GetAllProjects(false)
	if err != nil {
		return nil, err
	}

	var meetings []models.MeetingWithProject
	for _, proj := range projects {
		for _, m := range proj.Meetings {
			meetDate, _ := time.Parse("2006-01-02", m.MeetingDate)
			if (meetDate.Equal(startDate) || meetDate.After(startDate)) && meetDate.Before(endDate) {
				meetings = append(meetings, models.MeetingWithProject{
					Meeting:     m,
					ProjectName: proj.Name,
				})
			}
		}
	}

	return meetings, nil
}

// SearchMeetings searches for meetings by title, notes, or attendees
func (s *Service) SearchMeetings(query string) ([]models.MeetingWithProject, error) {
	return s.store.SearchAllMeetings(query)
}
