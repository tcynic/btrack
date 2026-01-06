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
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
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
