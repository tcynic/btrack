package notes

import (
	"strings"

	"btrack/internal/models"
)

// CreateTask creates a new task
func (s *Service) CreateTask(input models.CreateTaskInput) (*models.Task, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	task := &models.Task{
		ProjectID:   input.ProjectID,
		SourceType:  input.SourceType,
		SourceID:    input.SourceID,
		Title:       input.Title,
		Description: input.Description,
		Status:      models.TaskStatusPending,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
	}

	if err := s.store.AddTask(input.ProjectID, task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetTask returns a single task by ID
func (s *Service) GetTask(id int64) (*models.Task, error) {
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, t := range proj.Tasks {
			if t.ID == id {
				return &t, nil
			}
		}
	}
	return nil, models.NotFound("task")
}

// GetTasksByProject returns all tasks for a project
func (s *Service) GetTasksByProject(projectID int64) ([]models.Task, error) {
	return s.store.GetTasks(projectID)
}

// GetTasksBySource returns tasks linked to a specific source (meeting or note)
func (s *Service) GetTasksBySource(sourceType string, sourceID int64) ([]models.Task, error) {
	projects, _ := s.store.GetAllProjects(false)
	var tasks []models.Task
	for _, proj := range projects {
		for _, t := range proj.Tasks {
			if t.SourceType == sourceType && t.SourceID == sourceID {
				tasks = append(tasks, t)
			}
		}
	}
	return tasks, nil
}

// GetAllTasks returns all tasks across projects with optional filters
func (s *Service) GetAllTasks(statusFilter string, projectIDFilter int64) ([]models.TaskWithContext, error) {
	projects, err := s.store.GetAllProjects(false)
	if err != nil {
		return nil, err
	}

	var tasks []models.TaskWithContext
	for _, proj := range projects {
		for _, t := range proj.Tasks {
			// Apply filters
			if statusFilter != "" && t.Status != statusFilter {
				continue
			}
			if projectIDFilter > 0 && t.ProjectID != projectIDFilter {
				continue
			}

			// Get source name based on type
			var sourceName string
			if t.SourceType == "meeting" {
				for _, m := range proj.Meetings {
					if m.ID == t.SourceID {
						sourceName = m.Title
						break
					}
				}
			} else if t.SourceType == "note" {
				for _, n := range proj.Notes {
					if n.ID == t.SourceID {
						sourceName = n.Title
						break
					}
				}
			}

			tasks = append(tasks, models.TaskWithContext{
				Task:        t,
				ProjectName: proj.Name,
				SourceName:  sourceName,
			})
		}
	}

	// Sort: in_progress > pending > completed, then by priority desc, then by due date
	// Simple sort implementation
	for i := 0; i < len(tasks); i++ {
		for j := i + 1; j < len(tasks); j++ {
			if taskSortLess(tasks[j], tasks[i]) {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			}
		}
	}

	return tasks, nil
}

func taskSortLess(a, b models.TaskWithContext) bool {
	// Status priority
	statusOrder := map[string]int{
		"in_progress": 1,
		"pending":     2,
		"completed":   3,
	}
	aStatus := statusOrder[a.Status]
	bStatus := statusOrder[b.Status]
	if aStatus != bStatus {
		return aStatus < bStatus
	}
	// Priority (high > medium > low)
	priorityOrder := map[string]int{"high": 3, "medium": 2, "low": 1}
	aPrio := priorityOrder[a.Priority]
	bPrio := priorityOrder[b.Priority]
	if aPrio != bPrio {
		return aPrio > bPrio
	}
	// Due date
	return a.DueDate < b.DueDate
}

// UpdateTask updates an existing task
func (s *Service) UpdateTask(input models.UpdateTaskInput) (*models.Task, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	// Find the project this task belongs to
	var projectID int64
	var sourceType string
	var sourceID int64
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, t := range proj.Tasks {
			if t.ID == input.ID {
				projectID = proj.ID
				sourceType = t.SourceType
				sourceID = t.SourceID
				break
			}
		}
		if projectID > 0 {
			break
		}
	}

	if projectID == 0 {
		return nil, models.NotFound("task")
	}

	task := &models.Task{
		ID:          input.ID,
		ProjectID:   projectID,
		SourceType:  sourceType,
		SourceID:    sourceID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
	}

	if err := s.store.UpdateTask(projectID, task); err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateTaskStatus updates only the status of a task
func (s *Service) UpdateTaskStatus(id int64, status string) (*models.Task, error) {
	if !models.IsValidTaskStatus(status) {
		return nil, models.ValidationError("status", "invalid status value")
	}

	// Get the task first
	task, err := s.GetTask(id)
	if err != nil {
		return nil, err
	}

	// Update status
	task.Status = status

	// Find project and update
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, t := range proj.Tasks {
			if t.ID == id {
				if err := s.store.UpdateTask(proj.ID, task); err != nil {
					return nil, err
				}
				return task, nil
			}
		}
	}

	return nil, models.NotFound("task")
}

// DeleteTask removes a task
func (s *Service) DeleteTask(id int64) error {
	projects, _ := s.store.GetAllProjects(false)
	for _, proj := range projects {
		for _, t := range proj.Tasks {
			if t.ID == id {
				return s.store.DeleteTask(proj.ID, id)
			}
		}
	}
	return models.NotFound("task")
}

// SearchTasks searches for tasks by title or description
func (s *Service) SearchTasks(query string) ([]models.TaskWithContext, error) {
	return s.store.SearchAllTasks(query)
}
