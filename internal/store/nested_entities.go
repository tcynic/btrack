package store

import (
	"time"

	"btrack/internal/models"
)

// ===== MEETINGS =====

// AddMeeting adds a meeting to a project
func (s *Store) AddMeeting(projectID int64, meeting models.Meeting) (*models.Meeting, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	now := time.Now().Format(time.RFC3339)
	meeting.ID = s.data.NextIDs["meeting"]
	meeting.ProjectID = projectID
	meeting.CreatedAt = now
	meeting.UpdatedAt = now

	s.data.NextIDs["meeting"]++
	s.data.Projects[projectIdx].Meetings = append(s.data.Projects[projectIdx].Meetings, meeting)

	if err := s.saveUnlocked(); err != nil {
		return nil, err
	}

	return &meeting, nil
}

// GetMeetings returns all meetings for a project
func (s *Store) GetMeetings(projectID int64) ([]models.Meeting, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	meetings := s.data.Projects[projectIdx].Meetings
	if meetings == nil {
		meetings = []models.Meeting{}
	}

	return meetings, nil
}

// GetMeeting returns a specific meeting
func (s *Store) GetMeeting(projectID, meetingID int64) (*models.Meeting, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Meetings {
		if s.data.Projects[projectIdx].Meetings[i].ID == meetingID {
			return &s.data.Projects[projectIdx].Meetings[i], nil
		}
	}

	return nil, models.NotFound("meeting")
}

// UpdateMeeting updates a meeting
func (s *Store) UpdateMeeting(projectID int64, meeting models.Meeting) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Meetings {
		if s.data.Projects[projectIdx].Meetings[i].ID == meeting.ID {
			meeting.UpdatedAt = time.Now().Format(time.RFC3339)
			meeting.ProjectID = projectID
			// Preserve CreatedAt
			meeting.CreatedAt = s.data.Projects[projectIdx].Meetings[i].CreatedAt
			s.data.Projects[projectIdx].Meetings[i] = meeting
			return s.saveUnlocked()
		}
	}

	return models.NotFound("meeting")
}

// DeleteMeeting removes a meeting
func (s *Store) DeleteMeeting(projectID, meetingID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Meetings {
		if s.data.Projects[projectIdx].Meetings[i].ID == meetingID {
			meetings := s.data.Projects[projectIdx].Meetings
			s.data.Projects[projectIdx].Meetings = append(meetings[:i], meetings[i+1:]...)
			return s.saveUnlocked()
		}
	}

	return models.NotFound("meeting")
}

// ===== NOTES =====

// AddNote adds a note to a project
func (s *Store) AddNote(projectID int64, note models.Note) (*models.Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	now := time.Now().Format(time.RFC3339)
	note.ID = s.data.NextIDs["note"]
	note.ProjectID = projectID
	note.CreatedAt = now
	note.UpdatedAt = now

	s.data.NextIDs["note"]++
	s.data.Projects[projectIdx].Notes = append(s.data.Projects[projectIdx].Notes, note)

	if err := s.saveUnlocked(); err != nil {
		return nil, err
	}

	return &note, nil
}

// GetNotes returns all notes for a project
func (s *Store) GetNotes(projectID int64) ([]models.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	notes := s.data.Projects[projectIdx].Notes
	if notes == nil {
		notes = []models.Note{}
	}

	return notes, nil
}

// GetNote returns a specific note
func (s *Store) GetNote(projectID, noteID int64) (*models.Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Notes {
		if s.data.Projects[projectIdx].Notes[i].ID == noteID {
			return &s.data.Projects[projectIdx].Notes[i], nil
		}
	}

	return nil, models.NotFound("note")
}

// UpdateNote updates a note
func (s *Store) UpdateNote(projectID int64, note models.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Notes {
		if s.data.Projects[projectIdx].Notes[i].ID == note.ID {
			note.UpdatedAt = time.Now().Format(time.RFC3339)
			note.ProjectID = projectID
			// Preserve CreatedAt
			note.CreatedAt = s.data.Projects[projectIdx].Notes[i].CreatedAt
			s.data.Projects[projectIdx].Notes[i] = note
			return s.saveUnlocked()
		}
	}

	return models.NotFound("note")
}

// DeleteNote removes a note
func (s *Store) DeleteNote(projectID, noteID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Notes {
		if s.data.Projects[projectIdx].Notes[i].ID == noteID {
			notes := s.data.Projects[projectIdx].Notes
			s.data.Projects[projectIdx].Notes = append(notes[:i], notes[i+1:]...)
			return s.saveUnlocked()
		}
	}

	return models.NotFound("note")
}

// ===== GOALS =====

// AddGoal adds a goal to a project
func (s *Store) AddGoal(projectID int64, goal models.Goal) (*models.Goal, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	now := time.Now().Format(time.RFC3339)
	goal.ID = s.data.NextIDs["goal"]
	goal.ProjectID = projectID
	goal.CreatedAt = now
	goal.UpdatedAt = now

	s.data.NextIDs["goal"]++
	s.data.Projects[projectIdx].Goals = append(s.data.Projects[projectIdx].Goals, goal)

	if err := s.saveUnlocked(); err != nil {
		return nil, err
	}

	return &goal, nil
}

// GetGoals returns all goals for a project
func (s *Store) GetGoals(projectID int64) ([]models.Goal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	goals := s.data.Projects[projectIdx].Goals
	if goals == nil {
		goals = []models.Goal{}
	}

	return goals, nil
}

// GetGoal returns a specific goal
func (s *Store) GetGoal(projectID, goalID int64) (*models.Goal, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Goals {
		if s.data.Projects[projectIdx].Goals[i].ID == goalID {
			return &s.data.Projects[projectIdx].Goals[i], nil
		}
	}

	return nil, models.NotFound("goal")
}

// UpdateGoal updates a goal
func (s *Store) UpdateGoal(projectID int64, goal models.Goal) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Goals {
		if s.data.Projects[projectIdx].Goals[i].ID == goal.ID {
			goal.UpdatedAt = time.Now().Format(time.RFC3339)
			goal.ProjectID = projectID
			// Preserve CreatedAt
			goal.CreatedAt = s.data.Projects[projectIdx].Goals[i].CreatedAt
			s.data.Projects[projectIdx].Goals[i] = goal
			return s.saveUnlocked()
		}
	}

	return models.NotFound("goal")
}

// UpdateGoalStatus updates only the status of a goal
func (s *Store) UpdateGoalStatus(projectID, goalID int64, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Goals {
		if s.data.Projects[projectIdx].Goals[i].ID == goalID {
			s.data.Projects[projectIdx].Goals[i].Status = status
			s.data.Projects[projectIdx].Goals[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return s.saveUnlocked()
		}
	}

	return models.NotFound("goal")
}

// DeleteGoal removes a goal
func (s *Store) DeleteGoal(projectID, goalID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Goals {
		if s.data.Projects[projectIdx].Goals[i].ID == goalID {
			goals := s.data.Projects[projectIdx].Goals
			s.data.Projects[projectIdx].Goals = append(goals[:i], goals[i+1:]...)
			return s.saveUnlocked()
		}
	}

	return models.NotFound("goal")
}

// ===== TASKS =====

// AddTask adds a task to a project
func (s *Store) AddTask(projectID int64, task models.Task) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	now := time.Now().Format(time.RFC3339)
	task.ID = s.data.NextIDs["task"]
	task.ProjectID = projectID
	task.CreatedAt = now
	task.UpdatedAt = now

	s.data.NextIDs["task"]++
	s.data.Projects[projectIdx].Tasks = append(s.data.Projects[projectIdx].Tasks, task)

	if err := s.saveUnlocked(); err != nil {
		return nil, err
	}

	return &task, nil
}

// GetTasks returns all tasks for a project
func (s *Store) GetTasks(projectID int64) ([]models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	tasks := s.data.Projects[projectIdx].Tasks
	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

// GetTask returns a specific task
func (s *Store) GetTask(projectID, taskID int64) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Tasks {
		if s.data.Projects[projectIdx].Tasks[i].ID == taskID {
			return &s.data.Projects[projectIdx].Tasks[i], nil
		}
	}

	return nil, models.NotFound("task")
}

// GetTasksBySource returns tasks linked to a specific source
func (s *Store) GetTasksBySource(projectID int64, sourceType string, sourceID int64) ([]models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return nil, models.NotFound("project")
	}

	var result []models.Task
	for _, task := range s.data.Projects[projectIdx].Tasks {
		if task.SourceType == sourceType && task.SourceID != nil && *task.SourceID == sourceID {
			result = append(result, task)
		}
	}

	if result == nil {
		result = []models.Task{}
	}

	return result, nil
}

// UpdateTask updates a task
func (s *Store) UpdateTask(projectID int64, task models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Tasks {
		if s.data.Projects[projectIdx].Tasks[i].ID == task.ID {
			task.UpdatedAt = time.Now().Format(time.RFC3339)
			task.ProjectID = projectID
			// Preserve CreatedAt
			task.CreatedAt = s.data.Projects[projectIdx].Tasks[i].CreatedAt
			s.data.Projects[projectIdx].Tasks[i] = task
			return s.saveUnlocked()
		}
	}

	return models.NotFound("task")
}

// UpdateTaskStatus updates only the status of a task
func (s *Store) UpdateTaskStatus(projectID, taskID int64, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Tasks {
		if s.data.Projects[projectIdx].Tasks[i].ID == taskID {
			s.data.Projects[projectIdx].Tasks[i].Status = status
			s.data.Projects[projectIdx].Tasks[i].UpdatedAt = time.Now().Format(time.RFC3339)
			return s.saveUnlocked()
		}
	}

	return models.NotFound("task")
}

// DeleteTask removes a task
func (s *Store) DeleteTask(projectID, taskID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	for i := range s.data.Projects[projectIdx].Tasks {
		if s.data.Projects[projectIdx].Tasks[i].ID == taskID {
			tasks := s.data.Projects[projectIdx].Tasks
			s.data.Projects[projectIdx].Tasks = append(tasks[:i], tasks[i+1:]...)
			return s.saveUnlocked()
		}
	}

	return models.NotFound("task")
}

// DeleteTasksBySource removes tasks linked to a specific source
func (s *Store) DeleteTasksBySource(projectID int64, sourceType string, sourceID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectIdx := s.findProjectIndex(projectID)
	if projectIdx == -1 {
		return models.NotFound("project")
	}

	var kept []models.Task
	for _, task := range s.data.Projects[projectIdx].Tasks {
		// Keep tasks that don't match the source
		if task.SourceType != sourceType || task.SourceID == nil || *task.SourceID != sourceID {
			kept = append(kept, task)
		}
	}

	s.data.Projects[projectIdx].Tasks = kept
	return s.saveUnlocked()
}
