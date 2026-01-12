package database

// Project queries
const (
	InsertProject = `
		INSERT INTO projects (name, total_sold_hours, specialist_hours, start_date, end_date, is_active)
		VALUES (?, ?, ?, ?, ?, 1)
	`

	SelectAllProjects = `
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`

	SelectActiveProjects = `
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, created_at, updated_at
		FROM projects
		WHERE is_active = 1
		ORDER BY created_at DESC
	`

	SelectProjectByID = `
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, is_active, created_at, updated_at
		FROM projects
		WHERE id = ?
	`

	UpdateProject = `
		UPDATE projects
		SET name = ?, total_sold_hours = ?, specialist_hours = ?, start_date = ?, end_date = ?, is_active = ?, updated_at = datetime('now')
		WHERE id = ?
	`

	DeleteProject = `
		DELETE FROM projects WHERE id = ?
	`

	SoftDeleteProject = `
		UPDATE projects SET is_active = 0, updated_at = datetime('now') WHERE id = ?
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
)
