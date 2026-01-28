package store

import (
	"strings"

	"btrack/internal/models"
)

// SearchAllMeetings searches meetings across all projects
func (s *Store) SearchAllMeetings(query string) ([]models.MeetingWithProject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(query)
	var result []models.MeetingWithProject

	for _, project := range s.data.Projects {
		for _, meeting := range project.Meetings {
			// Search in title, notes, and attendees
			if strings.Contains(strings.ToLower(meeting.Title), query) ||
				strings.Contains(strings.ToLower(meeting.Notes), query) ||
				strings.Contains(strings.ToLower(meeting.Attendees), query) {
				result = append(result, models.MeetingWithProject{
					Meeting:     meeting,
					ProjectName: project.Name,
				})
			}
		}
	}

	if result == nil {
		result = []models.MeetingWithProject{}
	}

	return result, nil
}

// SearchAllNotes searches notes across all projects
func (s *Store) SearchAllNotes(query string) ([]models.NoteWithProject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(query)
	var result []models.NoteWithProject

	for _, project := range s.data.Projects {
		for _, note := range project.Notes {
			// Search in title and content
			if strings.Contains(strings.ToLower(note.Title), query) ||
				strings.Contains(strings.ToLower(note.Content), query) {
				result = append(result, models.NoteWithProject{
					Note:        note,
					ProjectName: project.Name,
				})
			}
		}
	}

	if result == nil {
		result = []models.NoteWithProject{}
	}

	return result, nil
}

// SearchAllGoals searches goals across all projects
func (s *Store) SearchAllGoals(query string) ([]models.GoalWithProject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(query)
	var result []models.GoalWithProject

	for _, project := range s.data.Projects {
		for _, goal := range project.Goals {
			// Search in title and description
			if strings.Contains(strings.ToLower(goal.Title), query) ||
				strings.Contains(strings.ToLower(goal.Description), query) {
				result = append(result, models.GoalWithProject{
					Goal:        goal,
					ProjectName: project.Name,
				})
			}
		}
	}

	if result == nil {
		result = []models.GoalWithProject{}
	}

	return result, nil
}

// SearchAllTasks searches tasks across all projects
func (s *Store) SearchAllTasks(query string) ([]models.TaskWithContext, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(query)
	var result []models.TaskWithContext

	for _, project := range s.data.Projects {
		for _, task := range project.Tasks {
			// Search in title and description
			if strings.Contains(strings.ToLower(task.Title), query) ||
				strings.Contains(strings.ToLower(task.Description), query) {
				
				// Find source title if applicable
				sourceTitle := ""
				if task.SourceType == "meeting" && task.SourceID != nil {
					for _, meeting := range project.Meetings {
						if meeting.ID == *task.SourceID {
							sourceTitle = meeting.Title
							break
						}
					}
				} else if task.SourceType == "note" && task.SourceID != nil {
					for _, note := range project.Notes {
						if note.ID == *task.SourceID {
							sourceTitle = note.Title
							break
						}
					}
				}

				result = append(result, models.TaskWithContext{
					Task:        task,
					ProjectName: project.Name,
					SourceTitle: sourceTitle,
				})
			}
		}
	}

	if result == nil {
		result = []models.TaskWithContext{}
	}

	return result, nil
}

// GetAllTasks returns all tasks across projects with optional filtering
func (s *Store) GetAllTasks(status, projectID string) ([]models.TaskWithContext, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.TaskWithContext

	for _, project := range s.data.Projects {
		// Filter by project if specified
		if projectID != "" && projectID != string(rune(project.ID)) {
			continue
		}

		for _, task := range project.Tasks {
			// Filter by status if specified
			if status != "" && task.Status != status {
				continue
			}

			// Find source title if applicable
			sourceTitle := ""
			if task.SourceType == "meeting" && task.SourceID != nil {
				for _, meeting := range project.Meetings {
					if meeting.ID == *task.SourceID {
						sourceTitle = meeting.Title
						break
					}
				}
			} else if task.SourceType == "note" && task.SourceID != nil {
				for _, note := range project.Notes {
					if note.ID == *task.SourceID {
						sourceTitle = note.Title
						break
					}
				}
			}

			result = append(result, models.TaskWithContext{
				Task:        task,
				ProjectName: project.Name,
				SourceTitle: sourceTitle,
			})
		}
	}

	if result == nil {
		result = []models.TaskWithContext{}
	}

	return result, nil
}

// GetMeetingsByDate returns all meetings for a specific date across all projects
func (s *Store) GetMeetingsByDate(date string) ([]models.MeetingWithProject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.MeetingWithProject

	for _, project := range s.data.Projects {
		if !project.IsActive {
			continue
		}

		for _, meeting := range project.Meetings {
			if meeting.MeetingDate == date {
				result = append(result, models.MeetingWithProject{
					Meeting:     meeting,
					ProjectName: project.Name,
				})
			}
		}
	}

	if result == nil {
		result = []models.MeetingWithProject{}
	}

	return result, nil
}

// GetMeetingsByWeek returns all meetings in a date range across all projects
func (s *Store) GetMeetingsByWeek(startDate, endDate string) ([]models.MeetingWithProject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.MeetingWithProject

	for _, project := range s.data.Projects {
		if !project.IsActive {
			continue
		}

		for _, meeting := range project.Meetings {
			if meeting.MeetingDate >= startDate && meeting.MeetingDate < endDate {
				result = append(result, models.MeetingWithProject{
					Meeting:     meeting,
					ProjectName: project.Name,
				})
			}
		}
	}

	if result == nil {
		result = []models.MeetingWithProject{}
	}

	return result, nil
}
