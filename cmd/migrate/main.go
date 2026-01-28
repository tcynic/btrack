package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"btrack/internal/models"
	"btrack/internal/store"
)

func main() {
	// Parse flags
	dbPath := flag.String("db", "", "Path to SQLite database file")
	jsonPath := flag.String("json", "", "Path to output JSON file")
	dryRun := flag.Bool("dry-run", false, "Preview migration without writing")
	flag.Parse()

	// Set default paths if not provided
	if *dbPath == "" {
		homeDir, _ := os.UserHomeDir()
		*dbPath = filepath.Join(homeDir, "Library", "Application Support", "btrack", "btrack.db")
	}
	if *jsonPath == "" {
		homeDir, _ := os.UserHomeDir()
		*jsonPath = filepath.Join(homeDir, "Library", "Application Support", "btrack", "btrack-data.json")
	}

	log.Printf("Migration Tool: SQLite → JSON")
	log.Printf("Source DB: %s", *dbPath)
	log.Printf("Target JSON: %s", *jsonPath)
	log.Printf("Dry run: %v", *dryRun)
	log.Println()

	// Open SQLite database
	db, err := sql.Open("sqlite", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create migrator
	migrator := &Migrator{db: db}

	// Perform migration
	data, err := migrator.Migrate()
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Print summary
	log.Printf("\n✓ Migration completed successfully!")
	log.Printf("  Projects: %d", len(data.Projects))
	totalEntries := 0
	totalMeetings := 0
	totalNotes := 0
	totalGoals := 0
	totalTasks := 0
	for _, proj := range data.Projects {
		totalEntries += len(proj.WeeklyEntries)
		totalMeetings += len(proj.Meetings)
		totalNotes += len(proj.Notes)
		totalGoals += len(proj.Goals)
		totalTasks += len(proj.Tasks)
	}
	log.Printf("  Weekly Entries: %d", totalEntries)
	log.Printf("  Meetings: %d", totalMeetings)
	log.Printf("  Notes: %d", totalNotes)
	log.Printf("  Goals: %d", totalGoals)
	log.Printf("  Tasks: %d", totalTasks)
	log.Printf("  Templates: %d", len(data.Templates))

	if *dryRun {
		log.Println("\n[DRY RUN] Migration preview completed. No files written.")
		return
	}

	// Backup existing JSON file if it exists
	if _, err := os.Stat(*jsonPath); err == nil {
		backupPath := *jsonPath + ".backup-" + time.Now().Format("2006-01-02-150405")
		if err := os.Rename(*jsonPath, backupPath); err != nil {
			log.Fatalf("Failed to backup existing JSON: %v", err)
		}
		log.Printf("\n✓ Backed up existing JSON to: %s", backupPath)
	}

	// Write JSON file
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	if err := os.WriteFile(*jsonPath, jsonData, 0644); err != nil {
		log.Fatalf("Failed to write JSON file: %v", err)
	}

	log.Printf("✓ Wrote data to: %s", *jsonPath)
	log.Println("\n🎉 Migration complete! You can now use the new JSON-based application.")
}

type Migrator struct {
	db *sql.DB
}

func (m *Migrator) Migrate() (*store.Data, error) {
	data := &store.Data{
		SchemaVersion: 1,
		NextIDs: map[string]int64{
			"project":      1,
			"weekly_entry": 1,
			"meeting":      1,
			"note":         1,
			"goal":         1,
			"task":         1,
			"template":     1,
		},
		Projects:  []store.ProjectWithNested{},
		Templates: []store.Template{},
	}

	log.Println("Migrating projects...")
	projects, err := m.migrateProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate projects: %w", err)
	}
	data.Projects = projects

	log.Println("Migrating templates...")
	templates, err := m.migrateTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate templates: %w", err)
	}
	data.Templates = templates

	// Calculate next IDs based on existing data
	m.calculateNextIDs(data)

	return data, nil
}

func (m *Migrator) migrateProjects() ([]store.ProjectWithNested, error) {
	rows, err := m.db.Query(`
		SELECT id, name, total_sold_hours, specialist_hours, start_date, end_date, 
		       is_active, is_persistent, created_at, updated_at
		FROM projects
		WHERE deleted_at IS NULL
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []store.ProjectWithNested
	for rows.Next() {
		var proj store.ProjectWithNested
		err := rows.Scan(
			&proj.ID, &proj.Name, &proj.TotalSoldHours, &proj.SpecialistHours,
			&proj.StartDate, &proj.EndDate, &proj.IsActive, &proj.IsPersistent,
			&proj.CreatedAt, &proj.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Migrate nested entities for this project
		proj.WeeklyEntries, err = m.migrateWeeklyEntries(proj.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to migrate weekly entries for project %d: %w", proj.ID, err)
		}

		proj.Meetings, err = m.migrateMeetings(proj.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to migrate meetings for project %d: %w", proj.ID, err)
		}

		proj.Notes, err = m.migrateNotes(proj.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to migrate notes for project %d: %w", proj.ID, err)
		}

		proj.Goals, err = m.migrateGoals(proj.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to migrate goals for project %d: %w", proj.ID, err)
		}

		proj.Tasks, err = m.migrateTasks(proj.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to migrate tasks for project %d: %w", proj.ID, err)
		}

		projects = append(projects, proj)
		log.Printf("  ✓ Migrated project: %s (ID: %d)", proj.Name, proj.ID)
	}

	return projects, nil
}

func (m *Migrator) migrateWeeklyEntries(projectID int64) ([]models.WeeklyEntry, error) {
	rows, err := m.db.Query(`
		SELECT id, project_id, week_start_date, week_number, planned_hours, actual_hours, 
		       created_at, updated_at
		FROM weekly_entries
		WHERE project_id = ?
		ORDER BY week_start_date
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.WeeklyEntry
	for rows.Next() {
		var entry models.WeeklyEntry
		var weekNum int64

		err := rows.Scan(&entry.ID, &entry.ProjectID, &entry.WeekStartDate, &weekNum, 
			&entry.PlannedHours, &entry.ActualHours, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}
		entry.WeekNumber = int(weekNum)

		entries = append(entries, entry)
	}

	if entries == nil {
		entries = []models.WeeklyEntry{}
	}
	return entries, nil
}

func (m *Migrator) migrateMeetings(projectID int64) ([]models.Meeting, error) {
	rows, err := m.db.Query(`
		SELECT id, project_id, title, meeting_date, duration_minutes, attendees, notes,
		       created_at, updated_at
		FROM project_meetings
		WHERE project_id = ?
		ORDER BY meeting_date DESC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var meeting models.Meeting
		var duration int64
		var attendees, notes *string

		err := rows.Scan(&meeting.ID, &meeting.ProjectID, &meeting.Title, &meeting.MeetingDate, 
			&duration, &attendees, &notes, &meeting.CreatedAt, &meeting.UpdatedAt)
		if err != nil {
			return nil, err
		}
		meeting.DurationMinutes = int(duration)
		meeting.Attendees = nullableString(attendees)
		meeting.Notes = nullableString(notes)

		meetings = append(meetings, meeting)
	}

	if meetings == nil {
		meetings = []models.Meeting{}
	}
	return meetings, nil
}

func (m *Migrator) migrateNotes(projectID int64) ([]models.Note, error) {
	rows, err := m.db.Query(`
		SELECT id, project_id, title, content, created_at, updated_at
		FROM project_notes
		WHERE project_id = ?
		ORDER BY created_at DESC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		var content *string

		err := rows.Scan(&note.ID, &note.ProjectID, &note.Title, &content, 
			&note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		note.Content = nullableString(content)

		notes = append(notes, note)
	}

	if notes == nil {
		notes = []models.Note{}
	}
	return notes, nil
}

func (m *Migrator) migrateGoals(projectID int64) ([]models.Goal, error) {
	rows, err := m.db.Query(`
		SELECT id, project_id, title, description, status, target_date, created_at, updated_at
		FROM project_goals
		WHERE project_id = ?
		ORDER BY created_at DESC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []models.Goal
	for rows.Next() {
		var goal models.Goal
		var description, targetDate *string

		err := rows.Scan(&goal.ID, &goal.ProjectID, &goal.Title, &description, &goal.Status, 
			&targetDate, &goal.CreatedAt, &goal.UpdatedAt)
		if err != nil {
			return nil, err
		}
		goal.Description = nullableString(description)
		goal.TargetDate = nullableString(targetDate)

		goals = append(goals, goal)
	}

	if goals == nil {
		goals = []models.Goal{}
	}
	return goals, nil
}

func (m *Migrator) migrateTasks(projectID int64) ([]models.Task, error) {
	rows, err := m.db.Query(`
		SELECT id, project_id, source_type, source_id, title, description, status, 
		       priority, due_date, created_at, updated_at
		FROM tasks
		WHERE project_id = ?
		ORDER BY created_at DESC
	`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var description, dueDate *string

		err := rows.Scan(&task.ID, &task.ProjectID, &task.SourceType, &task.SourceID, &task.Title, 
			&description, &task.Status, &task.Priority, &dueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		task.Description = nullableString(description)
		task.DueDate = nullableString(dueDate)

		tasks = append(tasks, task)
	}

	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}

func (m *Migrator) migrateTemplates() ([]store.Template, error) {
	rows, err := m.db.Query(`
		SELECT id, name, total_sold_hours, specialist_hours, created_at, updated_at
		FROM project_templates
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []store.Template
	for rows.Next() {
		var tmpl store.Template
		err := rows.Scan(&tmpl.ID, &tmpl.Name, &tmpl.TotalSoldHours, &tmpl.SpecialistHours,
			&tmpl.CreatedAt, &tmpl.UpdatedAt)
		if err != nil {
			return nil, err
		}
		templates = append(templates, tmpl)
	}

	if templates == nil {
		templates = []store.Template{}
	}
	return templates, nil
}

func (m *Migrator) calculateNextIDs(data *store.Data) {
	maxIDs := map[string]int64{
		"project":      0,
		"weekly_entry": 0,
		"meeting":      0,
		"note":         0,
		"goal":         0,
		"task":         0,
		"template":     0,
	}

	// Find max IDs from projects and nested entities
	for _, proj := range data.Projects {
		if proj.ID > maxIDs["project"] {
			maxIDs["project"] = proj.ID
		}
		for _, entry := range proj.WeeklyEntries {
			if entry.ID > maxIDs["weekly_entry"] {
				maxIDs["weekly_entry"] = entry.ID
			}
		}
		for _, meeting := range proj.Meetings {
			if meeting.ID > maxIDs["meeting"] {
				maxIDs["meeting"] = meeting.ID
			}
		}
		for _, note := range proj.Notes {
			if note.ID > maxIDs["note"] {
				maxIDs["note"] = note.ID
			}
		}
		for _, goal := range proj.Goals {
			if goal.ID > maxIDs["goal"] {
				maxIDs["goal"] = goal.ID
			}
		}
		for _, task := range proj.Tasks {
			if task.ID > maxIDs["task"] {
				maxIDs["task"] = task.ID
			}
		}
	}

	// Find max template ID
	for _, tmpl := range data.Templates {
		if tmpl.ID > maxIDs["template"] {
			maxIDs["template"] = tmpl.ID
		}
	}

	// Set next IDs (max + 1)
	for key, maxID := range maxIDs {
		data.NextIDs[key] = maxID + 1
	}
}

func nullableString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
