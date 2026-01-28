# SQLite to JSON Store Migration Guide

## Overview
This guide documents the migration from SQLite to an in-memory JSON-backed data store for btrack.

## ✅ Completed Work

### Store Package Implementation
Location: `internal/store/`

The complete store package has been implemented with 8 files (~1,770 lines):

1. **store.go** - Core functionality
   - `Load()` - Reads JSON file into memory
   - `Save()` - Atomic write (temp file + rename)
   - `NextID()` - Auto-increment ID generation

2. **types.go** - Data structures
   - `Data` - Top-level structure
   - `ProjectWithNested` - Projects with nested entities
   - `Template` - Project templates

3. **projects.go** (~260 lines)
   - Project CRUD operations
   - Search by name
   - Persistent project seeding
   - Project statistics

4. **weekly_entries.go** (~156 lines)
   - Weekly entry CRUD
   - Update actual hours
   - Delete future entries

5. **nested_entities.go** (~535 lines)
   - Meetings CRUD
   - Notes CRUD
   - Goals CRUD (with status updates)
   - Tasks CRUD (with source linking)

6. **search.go** (~250 lines)
   - Cross-project search for all entity types
   - Meeting queries by date/week
   - Task filtering by status/project

7. **aggregations.go** (~336 lines)
   - Dashboard weekly aggregates
   - Monthly trends
   - Variance reports
   - Capacity utilization
   - Goal statistics

8. **templates.go** (~74 lines)
   - Template CRUD operations

### Data Structure
All data stored in project-centric nested format:
- Projects contain: weekly_entries, meetings, notes, goals, tasks
- Templates are top-level (not tied to projects)
- Cascade deletes are automatic (delete project = delete all nested data)

## 🚧 Remaining Work

### Phase 3: Update Services

#### 3.1 Project Service
File: `internal/services/project/service.go`

**Changes:**
```go
// Replace constructor
func NewService(store *store.Store) *Service {
    return &Service{store: store}
}

// Update methods to use store
// GetByID now returns *ProjectWithNested with embedded weekly_entries
```

**Key Difference:** Weekly entries are now part of the project object, no separate query needed.

#### 3.2 Tracking Service  
Files: `service.go`, `dashboard.go`, `reports.go`, `weekly.go`

**SQL → Store Method Mappings:**
- `db.Query(SelectWeeklyAggregates)` → `store.GetWeeklyAggregates(start, end)`
- `db.Query(SelectMonthlyTrends)` → `store.GetMonthlyTrends(monthsBack)`
- `db.Query(SelectVarianceReport)` → `store.GetVarianceReport(start, end)`
- `db.Query(SelectCapacityUtilization)` → `store.GetCapacityUtilization(start, end)`
- `db.QueryRow(SelectActiveProjectCount)` → `store.GetActiveProjectCount()`
- `db.QueryRow(SelectAllGoalStats)` → `store.GetGoalStats()`

**Benefit:** Store returns structured types, no manual row scanning.

#### 3.3 Notes Service
Files: `service.go`, `meeting.go`, `note.go`, `goal.go`, `task.go`

**Pattern Change:** All nested entity methods now require `projectID` as first parameter.

Example:
```go
// OLD: repo.GetByProject(projectID)
// NEW: store.GetMeetings(projectID)

// OLD: repo.Create(meeting)  
// NEW: store.AddMeeting(projectID, meeting)
```

#### 3.4 System Service
Files: `service.go`, `template.go`, `backup.go`, `export.go`

**template.go:** Replace SQL with store template methods
**backup.go:** Update to backup/restore JSON file instead of SQLite
**export.go:** Use `store.GetAllProjects()` instead of SQL queries

### Phase 4: Update app.go

Replace database initialization:
```go
import "btrack/internal/store"

func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    
    // Get data file path
    configDir, _ := os.UserConfigDir()
    dataPath := filepath.Join(configDir, "btrack", "btrack-data.json")
    
    // Initialize store
    a.store = store.New(dataPath)
    if err := a.store.Load(); err != nil {
        log.Printf("Failed to load store: %v", err)
        return
    }
    
    // Seed persistent projects
    a.store.SeedPersistentProjects()
    
    // Initialize services
    a.projectService = project.NewService(a.store)
    a.trackingService = tracking.NewService(a.store)
    a.notesService = notes.NewService(a.store)
    a.systemService = system.NewService(a.projectService, a.store, ctx)
}

func (a *App) shutdown(ctx context.Context) {
    if a.store != nil {
        a.store.Save()
        log.Println("Store saved")
    }
}
```

### Phase 5: Migration Utility

Create `cmd/migrate/main.go`:
1. Open SQLite database
2. Query all tables
3. Build nested ProjectWithNested structures
4. Save to JSON via store

### Phase 6: Cleanup

1. Delete old packages:
   ```bash
   rm -rf internal/database
   rm -rf internal/repository
   ```

2. Clean dependencies:
   ```bash
   go mod tidy
   ```

3. Remove imports:
   - `"database/sql"`
   - `"btrack/internal/database"`
   - `"btrack/internal/repository"`

4. Update WARP.md documentation

## Benefits

- **~60% less code** - Eliminates database + repository layers (~1,400+ lines)
- **Simpler debugging** - Human-readable JSON file
- **Faster queries** - All data in memory
- **Easy backup** - Copy one file
- **No schema migrations** - Just code changes

## Trade-offs

- **Memory usage** - All data loaded at startup (negligible for typical usage)
- **Write amplification** - Entire file rewritten on changes (atomic writes via temp file)
- **Single process** - No concurrent access (already true for desktop app)

## Testing Checklist

After completing remaining phases:
- [ ] Create new project
- [ ] Add weekly entries  
- [ ] Add meetings/notes/goals/tasks
- [ ] Update actual hours
- [ ] Delete project (verify cascade)
- [ ] Search across projects
- [ ] View dashboard
- [ ] View reports
- [ ] Create/use template
- [ ] Backup/restore
- [ ] Export to CSV
- [ ] Restart app (verify persistence)
