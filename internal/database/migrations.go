package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a single database migration
type Migration struct {
	Version     int
	Description string
	SQL         string
}

// RunMigrations creates the database schema using versioned migrations
func RunMigrations(db *sql.DB) error {
	// Create schema_migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get current schema version
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	log.Printf("Current schema version: %d", currentVersion)

	// Define migrations in order
	migrations := []Migration{
		{1, "Create projects table", createProjectsTable},
		{2, "Create weekly entries table", createWeeklyEntriesTable},
		{3, "Create initial indexes", createIndexes},
		{4, "Create meetings table", createMeetingsTable},
		{5, "Create notes table", createNotesTable},
		{6, "Create goals table", createGoalsTable},
		{7, "Create additional indexes", createNewIndexes},
		{8, "Create templates table", createTemplatesTable},
		{9, "Create tasks table", createTasksTable},
		{10, "Create task indexes", createTaskIndexes},
		{11, "Add is_persistent column to projects", addPersistentColumnSQL},
		{12, "Add deleted_at column to projects", addDeletedAtColumnSQL},
	}

	// Apply migrations that haven't been applied yet
	for _, m := range migrations {
		if m.Version <= currentVersion {
			continue // Migration already applied
		}

		log.Printf("Applying migration %d: %s", m.Version, m.Description)
		if _, err := db.Exec(m.SQL); err != nil {
			return fmt.Errorf("migration %d failed: %w", m.Version, err)
		}

		// Record migration
		if err := recordMigration(db, m.Version, m.Description); err != nil {
			return fmt.Errorf("failed to record migration %d: %w", m.Version, err)
		}
	}

	log.Printf("Migrations complete. Schema version: %d", len(migrations))
	return nil
}

// createMigrationsTable creates the schema_migrations table
func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TEXT NOT NULL DEFAULT (datetime('now')),
			description TEXT
		)
	`)
	return err
}

// getCurrentVersion returns the highest applied migration version
func getCurrentVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	return version, err
}

// recordMigration records that a migration has been applied
func recordMigration(db *sql.DB, version int, description string) error {
	_, err := db.Exec(
		"INSERT INTO schema_migrations (version, description) VALUES (?, ?)",
		version, description,
	)
	return err
}

// GetSchemaVersion returns the current schema version (exported for health checks)
func GetSchemaVersion(db *sql.DB) (int, error) {
	return getCurrentVersion(db)
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

const createTemplatesTable = `
CREATE TABLE IF NOT EXISTS project_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    total_sold_hours INTEGER NOT NULL,
    specialist_hours INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

const createTasksTable = `
CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    source_type TEXT NOT NULL DEFAULT 'standalone',
    source_id INTEGER,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    priority TEXT NOT NULL DEFAULT 'medium',
    due_date TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);
`

const createTaskIndexes = `
CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_source ON tasks(source_type, source_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
`

// addPersistentColumnSQL adds is_persistent column (idempotent)
const addPersistentColumnSQL = `
ALTER TABLE projects ADD COLUMN is_persistent INTEGER NOT NULL DEFAULT 0;
`

// addDeletedAtColumnSQL adds deleted_at column for soft deletes (idempotent)
const addDeletedAtColumnSQL = `
ALTER TABLE projects ADD COLUMN deleted_at TEXT;
CREATE INDEX IF NOT EXISTS idx_projects_deleted_at ON projects(deleted_at);
`


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
