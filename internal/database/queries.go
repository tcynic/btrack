package database

// Project queries
const (
	InsertProject = `
		INSERT INTO projects (name, total_sold_hours, specialist_hours, start_date, end_date, is_active)
		VALUES (?, ?, ?, ?, ?, 1)
	`

	SelectAllProjects = `
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, is_persistent, created_at, updated_at
		FROM projects
		WHERE deleted_at IS NULL
		ORDER BY is_persistent DESC, created_at DESC
	`

	SelectActiveProjects = `
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, is_persistent, created_at, updated_at
		FROM projects
		WHERE is_active = 1 AND deleted_at IS NULL
		ORDER BY is_persistent DESC, created_at DESC
	`

	SelectProjectByID = `
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, is_persistent, created_at, updated_at
		FROM projects
		WHERE id = ? AND deleted_at IS NULL
	`

	UpdateProject = `
		UPDATE projects
		SET name = ?, total_sold_hours = ?, specialist_hours = ?, start_date = ?, end_date = ?, is_active = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	SoftDeleteProject = `
		UPDATE projects 
		SET deleted_at = datetime('now'), updated_at = datetime('now') 
		WHERE id = ? AND deleted_at IS NULL
	`

	RestoreProject = `
		UPDATE projects 
		SET deleted_at = NULL, updated_at = datetime('now') 
		WHERE id = ?
	`

	PermanentlyDeleteProject = `
		DELETE FROM projects WHERE id = ?
	`
)

// Weekly entry queries
const (
	InsertWeeklyEntry = `
		INSERT INTO weekly_entries (project_id, week_start_date, week_number, planned_hours, actual_hours)
		VALUES (?, ?, ?, ?, 0)
	`

	SelectWeeklyEntriesByProject = `
		SELECT id, project_id, week_start_date, week_number, planned_hours, actual_hours, created_at, updated_at
		FROM weekly_entries
		WHERE project_id = ?
		ORDER BY week_number ASC
	`

	SelectWeeklyEntriesByDateRange = `
		SELECT we.id, we.project_id, we.week_start_date, we.week_number, we.planned_hours, we.actual_hours, we.created_at, we.updated_at
		FROM weekly_entries we
		INNER JOIN projects p ON we.project_id = p.id
		WHERE p.is_active = 1 AND we.week_start_date >= ? AND we.week_start_date <= ?
		ORDER BY we.week_start_date ASC
	`

	SelectWeeklyEntriesByWeek = `
		SELECT we.id, we.project_id, we.week_start_date, we.week_number, we.planned_hours, we.actual_hours,
		       we.created_at, we.updated_at, p.name as project_name
		FROM weekly_entries we
		INNER JOIN projects p ON we.project_id = p.id
		WHERE p.is_active = 1 AND we.week_start_date = ?
		ORDER BY p.name ASC
	`

	UpdateActualHours = `
		UPDATE weekly_entries
		SET actual_hours = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	DeleteWeeklyEntriesByProject = `
		DELETE FROM weekly_entries WHERE project_id = ?
	`

	DeleteFutureWeeklyEntries = `
		DELETE FROM weekly_entries
		WHERE project_id = ? AND week_start_date >= ?
	`

	SelectPastWeeklyHours = `
		SELECT COALESCE(SUM(planned_hours), 0)
		FROM weekly_entries
		WHERE project_id = ? AND week_start_date < ?
	`

	SelectProjectStats = `
		SELECT
			COALESCE(SUM(planned_hours), 0) as total_planned,
			COALESCE(SUM(actual_hours), 0) as total_actual,
			COUNT(*) as total_weeks
		FROM weekly_entries
		WHERE project_id = ?
	`
)

// Dashboard queries
const (
	SelectWeeklyAggregates = `
		SELECT
			we.week_start_date,
			SUM(we.planned_hours) as total_planned,
			SUM(we.actual_hours) as total_actual,
			COUNT(DISTINCT we.project_id) as project_count
		FROM weekly_entries we
		INNER JOIN projects p ON we.project_id = p.id
		WHERE p.is_active = 1 AND we.week_start_date >= ? AND we.week_start_date <= ?
		GROUP BY we.week_start_date
		ORDER BY we.week_start_date ASC
	`

	SelectActiveProjectCount = `
		SELECT COUNT(*) FROM projects WHERE is_active = 1
	`

	SelectAllGoalStats = `
		SELECT
			COUNT(*) as total_goals,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_goals
		FROM project_goals g
		INNER JOIN projects p ON g.project_id = p.id
		WHERE p.is_active = 1
	`
)

// Meeting queries
const (
	InsertMeeting = `
		INSERT INTO project_meetings (project_id, title, meeting_date, duration_minutes, attendees, notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	SelectMeetingsByProject = `
		SELECT id, project_id, title, meeting_date, duration_minutes, attendees, notes, created_at, updated_at
		FROM project_meetings
		WHERE project_id = ?
		ORDER BY meeting_date DESC
	`

	SelectMeetingByID = `
		SELECT id, project_id, title, meeting_date, duration_minutes, attendees, notes, created_at, updated_at
		FROM project_meetings
		WHERE id = ?
	`

	UpdateMeeting = `
		UPDATE project_meetings
		SET title = ?, meeting_date = ?, duration_minutes = ?, attendees = ?, notes = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	DeleteMeeting = `
		DELETE FROM project_meetings WHERE id = ?
	`

	SelectMeetingsByDate = `
		SELECT m.id, m.project_id, m.title, m.meeting_date, m.duration_minutes,
		       m.attendees, m.notes, m.created_at, m.updated_at, p.name as project_name
		FROM project_meetings m
		INNER JOIN projects p ON m.project_id = p.id
		WHERE m.meeting_date = ?
		ORDER BY m.duration_minutes DESC, m.title ASC
	`

	SelectMeetingsByWeek = `
		SELECT m.id, m.project_id, m.title, m.meeting_date, m.duration_minutes,
		       m.attendees, m.notes, m.created_at, m.updated_at, p.name as project_name
		FROM project_meetings m
		INNER JOIN projects p ON m.project_id = p.id
		WHERE m.meeting_date >= ? AND m.meeting_date < ?
		ORDER BY m.meeting_date ASC, m.duration_minutes DESC
	`
)

// Note queries
const (
	InsertNote = `
		INSERT INTO project_notes (project_id, title, content)
		VALUES (?, ?, ?)
	`

	SelectNotesByProject = `
		SELECT id, project_id, title, content, created_at, updated_at
		FROM project_notes
		WHERE project_id = ?
		ORDER BY updated_at DESC
	`

	SelectNoteByID = `
		SELECT id, project_id, title, content, created_at, updated_at
		FROM project_notes
		WHERE id = ?
	`

	UpdateNote = `
		UPDATE project_notes
		SET title = ?, content = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	DeleteNote = `
		DELETE FROM project_notes WHERE id = ?
	`
)

// Goal queries
const (
	InsertGoal = `
		INSERT INTO project_goals (project_id, title, description, status, target_date)
		VALUES (?, ?, ?, 'pending', ?)
	`

	SelectGoalsByProject = `
		SELECT id, project_id, title, description, status, target_date, created_at, updated_at
		FROM project_goals
		WHERE project_id = ?
		ORDER BY CASE status
			WHEN 'in_progress' THEN 1
			WHEN 'pending' THEN 2
			WHEN 'completed' THEN 3
		END, created_at DESC
	`

	SelectGoalByID = `
		SELECT id, project_id, title, description, status, target_date, created_at, updated_at
		FROM project_goals
		WHERE id = ?
	`

	UpdateGoal = `
		UPDATE project_goals
		SET title = ?, description = ?, status = ?, target_date = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	UpdateGoalStatus = `
		UPDATE project_goals
		SET status = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	DeleteGoal = `
		DELETE FROM project_goals WHERE id = ?
	`

	SelectGoalStatsByProject = `
		SELECT status, COUNT(*) as count
		FROM project_goals
		WHERE project_id = ?
		GROUP BY status
	`
)

// Search queries
const (
	SearchProjects = `
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, is_persistent, created_at, updated_at
		FROM projects
		WHERE name LIKE '%' || ? || '%' AND deleted_at IS NULL
		ORDER BY is_active DESC, created_at DESC
		LIMIT 50
	`

	SearchMeetings = `
		SELECT m.id, m.project_id, m.title, m.meeting_date, m.duration_minutes,
		       m.attendees, m.notes, m.created_at, m.updated_at, p.name as project_name
		FROM project_meetings m
		INNER JOIN projects p ON m.project_id = p.id
		WHERE m.title LIKE '%' || ? || '%' 
		   OR m.notes LIKE '%' || ? || '%'
		   OR m.attendees LIKE '%' || ? || '%'
		ORDER BY m.meeting_date DESC
		LIMIT 50
	`

	SearchNotes = `
		SELECT n.id, n.project_id, n.title, n.content, n.created_at, n.updated_at, p.name as project_name
		FROM project_notes n
		INNER JOIN projects p ON n.project_id = p.id
		WHERE n.title LIKE '%' || ? || '%'
		   OR n.content LIKE '%' || ? || '%'
		ORDER BY n.updated_at DESC
		LIMIT 50
	`

	SearchGoals = `
		SELECT g.id, g.project_id, g.title, g.description, g.status, g.target_date, g.created_at, g.updated_at, p.name as project_name
		FROM project_goals g
		INNER JOIN projects p ON g.project_id = p.id
		WHERE g.title LIKE '%' || ? || '%'
		   OR g.description LIKE '%' || ? || '%'
		ORDER BY CASE g.status
			WHEN 'in_progress' THEN 1
			WHEN 'pending' THEN 2
			WHEN 'completed' THEN 3
		END, g.created_at DESC
		LIMIT 50
	`
)

// Report queries
const (
	SelectMonthlyTrends = `
		SELECT
			strftime('%Y-%m', we.week_start_date) as month,
			SUM(we.planned_hours) as total_planned,
			SUM(we.actual_hours) as total_actual,
			COUNT(DISTINCT we.project_id) as project_count
		FROM weekly_entries we
		INNER JOIN projects p ON we.project_id = p.id
		WHERE p.is_active = 1 AND we.week_start_date >= ?
		GROUP BY month
		ORDER BY month ASC
	`

	SelectVarianceReport = `
		SELECT
			p.id,
			p.name,
			COALESCE(SUM(we.planned_hours), 0) as total_planned,
			COALESCE(SUM(we.actual_hours), 0) as total_actual
		FROM projects p
		LEFT JOIN weekly_entries we ON p.id = we.project_id 
			AND we.week_start_date >= ? 
			AND we.week_start_date <= ?
		WHERE p.is_active = 1
		GROUP BY p.id, p.name
		HAVING total_planned > 0
		ORDER BY total_planned DESC
	`

	SelectCapacityUtilization = `
		SELECT
			we.week_start_date,
			SUM(we.planned_hours) as total_planned,
			SUM(we.actual_hours) as total_actual,
			COUNT(DISTINCT we.project_id) as project_count
		FROM weekly_entries we
		INNER JOIN projects p ON we.project_id = p.id
		WHERE p.is_active = 1 
			AND we.week_start_date >= ? 
			AND we.week_start_date <= ?
		GROUP BY we.week_start_date
		ORDER BY we.week_start_date ASC
	`
)

// Template queries
const (
	InsertTemplate = `
		INSERT INTO project_templates (name, total_sold_hours, specialist_hours)
		VALUES (?, ?, ?)
	`

	SelectAllTemplates = `
		SELECT id, name, total_sold_hours, specialist_hours, created_at, updated_at
		FROM project_templates
		ORDER BY created_at DESC
	`

	SelectTemplateByID = `
		SELECT id, name, total_sold_hours, specialist_hours, created_at, updated_at
		FROM project_templates
		WHERE id = ?
	`

	DeleteTemplate = `
		DELETE FROM project_templates WHERE id = ?
	`
)

// Task queries
const (
	InsertTask = `
		INSERT INTO tasks (project_id, source_type, source_id, title, description, status, priority, due_date)
		VALUES (?, ?, ?, ?, ?, 'pending', ?, ?)
	`

	SelectTaskByID = `
		SELECT id, project_id, source_type, source_id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`

	SelectTasksByProject = `
		SELECT id, project_id, source_type, source_id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE project_id = ?
		ORDER BY CASE status
			WHEN 'in_progress' THEN 1
			WHEN 'pending' THEN 2
			WHEN 'completed' THEN 3
			WHEN 'cancelled' THEN 4
		END, priority DESC, due_date ASC
	`

	SelectTasksBySource = `
		SELECT id, project_id, source_type, source_id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE source_type = ? AND source_id = ?
		ORDER BY created_at DESC
	`

	SelectAllTasks = `
		SELECT t.id, t.project_id, t.source_type, t.source_id, t.title, t.description, t.status, t.priority, t.due_date,
		       t.created_at, t.updated_at, p.name as project_name,
		       COALESCE(m.title, n.title, '') as source_title
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN project_meetings m ON t.source_type = 'meeting' AND t.source_id = m.id
		LEFT JOIN project_notes n ON t.source_type = 'note' AND t.source_id = n.id
		WHERE p.is_active = 1
		ORDER BY CASE t.status
			WHEN 'in_progress' THEN 1
			WHEN 'pending' THEN 2
			WHEN 'completed' THEN 3
			WHEN 'cancelled' THEN 4
		END, t.priority DESC, t.due_date ASC
	`

	SelectAllTasksFiltered = `
		SELECT t.id, t.project_id, t.source_type, t.source_id, t.title, t.description, t.status, t.priority, t.due_date,
		       t.created_at, t.updated_at, p.name as project_name,
		       COALESCE(m.title, n.title, '') as source_title
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN project_meetings m ON t.source_type = 'meeting' AND t.source_id = m.id
		LEFT JOIN project_notes n ON t.source_type = 'note' AND t.source_id = n.id
		WHERE p.is_active = 1
	`

	UpdateTask = `
		UPDATE tasks
		SET title = ?, description = ?, status = ?, priority = ?, due_date = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	UpdateTaskStatus = `
		UPDATE tasks
		SET status = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	DeleteTask = `
		DELETE FROM tasks WHERE id = ?
	`

	DeleteTasksBySource = `
		DELETE FROM tasks WHERE source_type = ? AND source_id = ?
	`

	SearchTasks = `
		SELECT t.id, t.project_id, t.source_type, t.source_id, t.title, t.description, t.status, t.priority, t.due_date,
		       t.created_at, t.updated_at, p.name as project_name,
		       COALESCE(m.title, n.title, '') as source_title
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN project_meetings m ON t.source_type = 'meeting' AND t.source_id = m.id
		LEFT JOIN project_notes n ON t.source_type = 'note' AND t.source_id = n.id
		WHERE (t.title LIKE '%' || ? || '%' OR t.description LIKE '%' || ? || '%')
		ORDER BY t.created_at DESC
		LIMIT 50
	`
)
