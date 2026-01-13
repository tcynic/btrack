package database

import (
	"database/sql"
	"fmt"
)

// RunMigrations creates the database schema
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createProjectsTable,
		createWeeklyEntriesTable,
		createIndexes,
		createMeetingsTable,
		createNotesTable,
		createGoalsTable,
		createNewIndexes,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	// Handle is_persistent column migration separately (idempotent)
	if err := addPersistentColumnIfNotExists(db); err != nil {
		return fmt.Errorf("failed to add is_persistent column: %w", err)
	}

	return nil
}

const createProjectsTable = `
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    total_sold_hours INTEGER NOT NULL,
    specialist_hours INTEGER NOT NULL DEFAULT 0,
    start_date TEXT NOT NULL,
    end_date TEXT NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const createWeeklyEntriesTable = `
CREATE TABLE IF NOT EXISTS weekly_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    week_start_date TEXT NOT NULL,
    week_number INTEGER NOT NULL,
    planned_hours INTEGER NOT NULL,
    actual_hours INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    UNIQUE(project_id, week_start_date)
);
`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_weekly_entries_project_id ON weekly_entries(project_id);
CREATE INDEX IF NOT EXISTS idx_weekly_entries_week_start ON weekly_entries(week_start_date);
CREATE INDEX IF NOT EXISTS idx_projects_is_active ON projects(is_active);
`

const createMeetingsTable = `
CREATE TABLE IF NOT EXISTS project_meetings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    meeting_date TEXT NOT NULL,
    duration_minutes INTEGER NOT NULL DEFAULT 60,
    attendees TEXT,
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
`

const createNotesTable = `
CREATE TABLE IF NOT EXISTS project_notes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
`

const createGoalsTable = `
CREATE TABLE IF NOT EXISTS project_goals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    target_date TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
`

const createNewIndexes = `
CREATE INDEX IF NOT EXISTS idx_meetings_project_id ON project_meetings(project_id);
CREATE INDEX IF NOT EXISTS idx_meetings_date ON project_meetings(meeting_date);
CREATE INDEX IF NOT EXISTS idx_notes_project_id ON project_notes(project_id);
CREATE INDEX IF NOT EXISTS idx_goals_project_id ON project_goals(project_id);
CREATE INDEX IF NOT EXISTS idx_goals_status ON project_goals(status);
`

// addPersistentColumnIfNotExists adds is_persistent column only if it doesn't exist
func addPersistentColumnIfNotExists(db *sql.DB) error {
	// Check if column exists by querying table info
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM pragma_table_info('projects') 
		WHERE name = 'is_persistent'
	`).Scan(&count)
	if err != nil {
		return err
	}

	// Column already exists, skip
	if count > 0 {
		return nil
	}

	// Add the column
	_, err = db.Exec(`ALTER TABLE projects ADD COLUMN is_persistent INTEGER NOT NULL DEFAULT 0`)
	return err
}

// SeedPersistentProjects creates the two persistent work projects if they don't exist
func SeedPersistentProjects(db *sql.DB) error {
	persistentProjects := []struct {
		Name string
	}{
		{Name: "Management"},
		{Name: "Internal Projects"},
	}

	for _, proj := range persistentProjects {
		// Check if project already exists
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM projects WHERE name = ? AND is_persistent = 1", proj.Name).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check for existing persistent project: %w", err)
		}

		// Skip if already exists
		if count > 0 {
			continue
		}

		// Insert persistent project with special values
		_, err = db.Exec(`
			INSERT INTO projects (name, total_sold_hours, specialist_hours, start_date, end_date, is_active, is_persistent)
			VALUES (?, 0, 0, '1900-01-01', '2099-12-31', 1, 1)
		`, proj.Name)
		if err != nil {
			return fmt.Errorf("failed to insert persistent project %s: %w", proj.Name, err)
		}
	}

	return nil
}
