# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

btrack (Bandwidth Tracker) is a desktop application for tracking project bandwidth allocation across weeks. Built with Wails v2, combining a Go backend with a React/TypeScript frontend, and uses a JSON file for local data persistence.

**Core Functionality**: Track "My Hours" (TotalSoldHours - SpecialistHours) across weekly entries using a frontloaded distribution algorithm. Users can plan project timelines, view dashboard aggregations, and update actual hours for past weeks.

## Prerequisites

- Go 1.24+
- Node.js/npm
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

## Key Commands

### Development
```bash
wails dev                    # Live development with hot reload (both backend and frontend)
cd frontend && npm install   # Install frontend dependencies
cd frontend && npm run dev   # Run Vite dev server independently
cd frontend && npm run build # Build frontend assets only
```

**Note**: During `wails dev`, the frontend is accessible at the dev server URL, and Go methods are available at http://localhost:34115 for browser-based development.

### Building
```bash
wails build                  # Build production binary for current platform
```

Creates a native application binary with embedded frontend assets.

### Data Management
```bash
# View data file during development (macOS)
cat ~/Library/Application\ Support/btrack/btrack-data.json
# Reset data (creates backup on next run)
rm ~/Library/Application\ Support/btrack/btrack-data.json
```

**Data Location**: `~/Library/Application Support/btrack/btrack-data.json` (macOS)

## Architecture

### Backend Structure (Go)

**Entry Point & Lifecycle**:
- **main.go**: Configures Wails app (1280x800 window), embeds `frontend/dist` assets via `//go:embed all:frontend/dist`, binds App instance to frontend
- **app.go**: Main App struct with Store reference; `startup(ctx)` receives runtime context and loads data from JSON, `shutdown(ctx)` saves data to JSON
- **Binding**: All public (capitalized) methods on App struct are automatically callable from JavaScript/TypeScript frontend

**Domain Logic** (all methods exposed to frontend via Wails bindings):
- **project.go**: CRUD operations for projects with automatic weekly entry generation/recalculation
  - `CreateProject`: Creates project and generates frontloaded weekly entries
  - `UpdateProject`: Updates project; recalculates future weeks if hours/dates change, preserves past weeks (skips recalc for persistent projects)
  - `GetAllProjects`, `GetProject`, `DeleteProject`, `ToggleProjectActive`
  - `SearchProjects`: Search projects by name
- **calculator.go**: Frontloaded distribution algorithm
  - `CalculateDistribution`: MyHours = TotalSoldHours - SpecialistHours, distributed across weeks with remainder frontloaded to earliest weeks
  - Helper functions: `calculateWeeksBetween`, `getMonday`, `getCurrentWeekMonday`
- **weekly_entry.go**: Weekly entry operations
  - `GetWeeklyEntries`: Returns entries with status (past/current/future week)
  - `UpdateActualHours`: Updates actual hours (only allowed for past/current weeks)
  - `CreateWeeklyEntry`: Manually create weekly entries (for persistent projects)
- **dashboard.go**: Aggregated statistics
  - `GetDashboardData`: Returns weekly aggregates across all projects
  - `GetDashboardSummary`: High-level stats (active projects, this/next week totals)
- **meeting.go**: CRUD operations for project meetings
  - `CreateMeeting`, `GetMeetings`, `GetMeeting`, `UpdateMeeting`, `DeleteMeeting`
  - `GetMeetingsByDate`: Returns all meetings for a specific date across projects
  - `GetMeetingsByWeek`: Returns all meetings for a 7-day period
  - `SearchMeetings`: Search meetings by title, notes, or attendees
- **note.go**: CRUD operations for project notes (markdown support in frontend)
  - `CreateNote`, `GetNotes`, `GetNote`, `UpdateNote`, `DeleteNote`
  - `SearchNotes`: Search notes by title or content (returns NoteWithProject)
- **goal.go**: CRUD operations for project goals with status tracking
  - `CreateGoal`, `GetGoals`, `GetGoal`, `UpdateGoal`, `DeleteGoal`
  - `UpdateGoalStatus`: Update only status field (pending/in_progress/completed/cancelled)
  - `GetGoalStats`: Returns goal statistics for a project (counts by status, completion rate)
  - `SearchGoals`: Search goals by title or description (returns GoalWithProject)
- **task.go**: Task management with linking capabilities
  - `CreateTask`, `GetTask`, `GetTasksByProject`, `GetTasksBySource`
  - `GetAllTasks`: Returns all tasks across projects with filters (status, project)
  - `UpdateTask`, `UpdateTaskStatus`, `DeleteTask`
  - `SearchTasks`: Search tasks by title or description (returns TaskWithContext)
  - Tasks can be standalone or linked to meetings/notes via source_type and source_id
  - Status values: pending, in_progress, completed
  - Priority values: low, medium, high
- **template.go**: Project template management
  - `CreateTemplate`: Save a project as a reusable template
  - `GetTemplates`, `GetTemplate`: Retrieve templates
  - `CreateProjectFromTemplate`: Create new project from template with custom dates
  - `DeleteTemplate`: Remove template
- **backup.go**: Database backup and restore
  - `GetBackupInfo`: Returns database path, size, last modified time
  - `CreateBackup`: Opens save dialog, copies database to user-selected location
  - `RestoreBackup`: Opens file dialog, restores database from backup (creates temp backup first)
- **export.go**: CSV export functionality
  - `ExportWeeklyReport`: Export weekly hours data for date range
  - `ExportProjectSummary`: Export single project data with all weekly entries
  - `ExportAllProjects`: Export summary data for all active projects
- **reports.go**: Analytics and reporting
  - `GetMonthlyTrends`: Returns monthly aggregated data (planned/actual/variance) for N months back
  - `GetVarianceReport`: Planned vs actual analysis by project for date range
  - `GetCapacityUtilization`: Weekly capacity utilization data with percentages

**Data Layer**:
- **internal/store/**: In-memory state with JSON persistence
  - `store.go`: Store struct with Load/Save methods and atomic write-to-temp-then-rename pattern
  - `types.go`: Data struct with project-centric nested structure (ProjectWithNested wraps Project + nested entities)
  - `projects.go`: Project CRUD operations (~260 lines)
  - `weekly_entries.go`: Weekly entry operations
  - `nested_entities.go`: CRUD for meetings, notes, goals, tasks (~535 lines)
  - `search.go`: Cross-project search operations
  - `aggregations.go`: Dashboard, reports, statistics (~336 lines)
  - `templates.go`: Template CRUD operations
  - All operations lock/unlock for thread-safety, save after mutations
- **internal/models/**: Domain models and validation
  - `project.go`: Project, ProjectWithStats, ProjectHealth, CreateProjectInput, UpdateProjectInput
  - `weekly_entry.go`: WeeklyEntry, WeeklyEntryWithStatus (includes IsPastWeek, Status calculations)
  - `meeting.go`: Meeting, MeetingWithProject, CreateMeetingInput, UpdateMeetingInput
  - `note.go`: Note, CreateNoteInput, UpdateNoteInput
  - `goal.go`: Goal, CreateGoalInput, UpdateGoalInput (with status constants: pending/in_progress/completed/cancelled)
  - `task.go`: Task, TaskWithContext, CreateTaskInput, UpdateTaskInput (with status: pending/in_progress/completed; priority: low/medium/high)
  - `search.go`: NoteWithProject, GoalWithProject, TaskWithContext (models for search results with parent project names)
  - `errors.go`: Custom error types for all domain objects
  - `nullable.go`: NullableString helper for handling optional string fields

### Frontend Structure (React/TypeScript + Vite)

- **Stack**: React 18, TypeScript, Vite, Tailwind CSS, Recharts for visualizations
- **Location**: `frontend/src/`
  - `App.tsx`: Main component
  - `components/`: UI components
  - `hooks/`: Custom React hooks
  - `context/`: React context providers
  - `types/`: TypeScript type definitions
  - `utils/`: Utility functions
- **Wails Bindings**: `frontend/wailsjs/` (auto-generated, DO NOT manually edit)
  - Import Go methods: `import {CreateProject} from '../wailsjs/go/main/App'`
  - Bindings regenerate automatically during `wails dev`

### Data Structure (JSON)

The application stores all data in a single JSON file with a project-centric nested structure:

```json
{
  "schema_version": 1,
  "next_ids": {"project": 11, "meeting": 26, ...},
  "projects": [
    {
      /* project fields */,
      "weekly_entries": [...],
      "meetings": [...],
      "notes": [...],
      "goals": [...],
      "tasks": [...]
    }
  ],
  "templates": [...]
}
```

**projects**:
- Stores client projects with total_sold_hours, specialist_hours, date range, is_active, is_persistent flags
- MyHours calculated as: TotalSoldHours - SpecialistHours
- Persistent projects ("Management", "Internal Projects"): special projects with dates 1900-2099, require manual weekly entry creation
- Each project contains nested arrays for all related entities

**weekly_entries** (nested under projects):
- One entry per project per week (Monday start dates)
- planned_hours: Calculated via frontloading algorithm (for non-persistent projects)
- actual_hours: User-entered for past/current weeks only
- Cascade delete: Removing project removes all nested weekly entries

**meetings** (nested under projects):
- Meeting records linked to projects: title, date, duration_minutes, attendees, notes
- Cascade delete with parent project

**notes** (nested under projects):
- Notes linked to projects: title, content (markdown)
- Cascade delete with parent project

**goals** (nested under projects):
- Goals linked to projects: title, description, status, target_date
- Status values: pending, in_progress, completed, cancelled
- Cascade delete with parent project

**tasks** (nested under projects):
- Tasks linked to projects: title, description, status, priority, due_date
- Can be standalone (source_type='standalone') or linked to meetings/notes (source_type='meeting'/'note' with source_id)
- Status values: pending, in_progress, completed
- Priority values: low, medium, high
- Cascade delete with parent project

**templates** (top-level):
- Saved project templates: name, total_sold_hours, specialist_hours
- Used to quickly create new projects with predefined hour allocations
- Independent of projects

## Code Patterns

### Backend (Go)

**Store Operations**: All CRUD operations go through the Store, which handles locking and persistence:
```go
// Store automatically locks, mutates, saves, and unlocks
project, err := s.CreateProject(input)
if err != nil {
  return nil, err
}
```

**Thread Safety**: Store uses sync.RWMutex for concurrent read/write safety. Read operations use RLock, write operations use Lock.

**Atomic Saves**: Store.Save() uses write-to-temp-then-rename pattern to ensure data integrity:
```go
// Write to .tmp file, then atomically rename
if err := s.saveUnlocked(); err != nil {
  return err
}
```

**Empty Slice Handling**: Methods return empty slices instead of nil to prevent frontend null errors:
```go
if result == nil {
  result = []Meeting{}
}
return result, nil
```

### Frontend (React/TypeScript)

**CRUD Hooks**: Use the `useCrud` factory for consistent async state management:
```typescript
const config = useMemo(() => ({
  loadFn: () => GetMeetings(projectId),
  createFn: CreateMeeting,
  updateFn: UpdateMeeting,
  deleteFn: DeleteMeeting,
  getId: (m: Meeting) => m.id,
}), [projectId])

const crud = useCrud<Meeting, CreateMeetingInput, UpdateMeetingInput>(config)
```

**Entity Lists**: Use `EntityList` component for consistent list rendering:
```typescript
<EntityList
  title="Meetings"
  items={meetings}
  isLoading={isLoading}
  emptyMessage="No meetings yet."
  addButtonLabel="Add Meeting"
  onAdd={() => setIsModalOpen(true)}
  renderItem={renderMeetingItem}
/>
```

## Development Workflow

### Adding Backend Methods
1. Add public method to App struct in `app.go` (or domain files: `project.go`, `calculator.go`, etc.)
2. Methods must be exported (capitalized) to be accessible from frontend
3. Run `wails dev` to regenerate bindings in `frontend/wailsjs/`
4. Import from `frontend/wailsjs/go/main/App` in React components

Example:
```go
// In app.go or domain file
func (a *App) MyNewMethod(param string) string {
    return "result: " + param
}
```

```typescript
// In React component
import {MyNewMethod} from '../wailsjs/go/main/App'
MyNewMethod("test").then(result => console.log(result))
```

### Migrating from SQLite (Legacy)
If you have existing SQLite data from a previous version:
```bash
# Build migration tool
go build -o migrate-tool ./cmd/migrate

# Preview migration (dry run)
./migrate-tool --dry-run

# Execute migration
./migrate-tool
```

See `cmd/migrate/README.md` for details.

### Frontend Development
- Edit components in `frontend/src/`
- Hot reload active during `wails dev`
- Use Wails bindings to call Go methods (they return Promises)
- Frontend build output goes to `frontend/dist` (embedded into binary)

## Critical Constraints

**Wails Bindings**: Files in `frontend/wailsjs/` are auto-generated by Wails. Manual edits will be overwritten.

**Go Method Visibility**: Only exported (capitalized) methods in the App struct are accessible from frontend.

**Weekly Entry Editing**: `UpdateActualHours` enforces that only past/current weeks can be edited (< next Monday).

**Project Updates**: When updating a project, the system preserves past/current week actual hours and only recalculates future weeks if hours/dates change. Persistent projects skip automatic recalculation entirely.

**Frontloading Algorithm**: Extra hours from MyHours % Weeks are distributed to the earliest weeks first (e.g., 10 hours / 3 weeks → 4, 3, 3).

**Persistent Projects**: "Management" and "Internal Projects" are special projects auto-seeded on first run with is_persistent=1. These don't use the frontloading algorithm and require manual weekly entry creation via `CreateWeeklyEntry`.

**Backup/Restore**: `CreateBackup` and `RestoreBackup` use Wails runtime dialogs for file selection. During restore, the current store is closed, a temporary backup is created (`.pre-restore`), the restore is performed, and the store is reloaded. If restore fails, the temp backup is used to rollback.

**Search Functionality**: Search methods (`SearchProjects`, `SearchNotes`, `SearchGoals`, `SearchMeetings`, `SearchTasks`) use case-insensitive string matching. Search results for notes, goals, meetings, and tasks include parent project names.

**CSV Export**: Export methods return CSV strings (not files). Frontend must handle saving the returned string to a file. Export includes headers and properly formatted data with variance calculations.

**Project Health**: Projects automatically calculate health status (on_track/at_risk/over_budget/completed) based on actual vs planned hours, end date, and current date. Persistent projects always show "on_track" status.

## Configuration Files

- **wails.json**: Wails project config (frontend commands, app metadata)
- **go.mod**: Go 1.24, Wails v2.11.0
- **frontend/package.json**: React 18, TypeScript, Vite, Tailwind, Recharts
- **frontend/vite.config.js**: Vite configuration
- **frontend/tailwind.config.js**: Tailwind CSS configuration
