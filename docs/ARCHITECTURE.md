# btrack Architecture Documentation

## Overview

btrack (Bandwidth Tracker) is a desktop application for tracking project bandwidth allocation across weeks. It uses **Wails v2** to combine a Go backend with a React/TypeScript frontend, with all data stored in a single JSON file for simplicity and portability.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  Desktop Application                     │
│                                                          │
│  ┌────────────────────────────────────────────────────┐ │
│  │         Frontend (React/TypeScript)                 │ │
│  │  • UI Components                                    │ │
│  │  • Custom Hooks (useCrud, useProjects, etc.)       │ │
│  │  • Context Providers                                │ │
│  └──────────────────┬─────────────────────────────────┘ │
│                     │                                    │
│  ┌──────────────────▼─────────────────────────────────┐ │
│  │         Wails Runtime                               │ │
│  │  • Auto-generated Bindings                         │ │
│  │  • Event System                                     │ │
│  │  • Native Dialogs                                   │ │
│  └──────────────────┬─────────────────────────────────┘ │
│                     │                                    │
│  ┌──────────────────▼─────────────────────────────────┐ │
│  │         Backend (Go)                                │ │
│  │  • App Struct (lifecycle & method bindings)        │ │
│  │  • Services (project, tracking, notes, system)     │ │
│  │  • Store (in-memory state + JSON persistence)      │ │
│  │  • Models (domain logic & validation)              │ │
│  └──────────────────┬─────────────────────────────────┘ │
│                     │                                    │
│  ┌──────────────────▼─────────────────────────────────┐ │
│  │    JSON File Storage (btrack-data.json)            │ │
│  │  • Project-centric nested structure                │ │
│  │  • Atomic writes (temp file + rename)              │ │
│  │  • Human-readable format                            │ │
│  └────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## Technology Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Frontend | React 18, TypeScript, Vite | UI rendering and state management |
| Styling | Tailwind CSS | Utility-first CSS framework |
| Charts | Recharts | Data visualization |
| Desktop Framework | Wails v2 | Desktop app packaging and Go-JS bridge |
| Backend | Go 1.24+ | Business logic and data operations |
| Storage | JSON File | Local data persistence (in-memory state) |

## Backend Architecture

### Application Lifecycle

```
┌──────────┐         ┌──────────────┐         ┌─────────┐         ┌──────────┐
│ main.go  │────────>│ Wails Runtime│────────>│ App     │────────>│ Store    │
└──────────┘         └──────────────┘         └─────────┘         └──────────┘
                            │                       │                    │
                            │ wails.Run()           │ startup(ctx)       │
                            │                       │                    │
                            │                       │──> store.Load() ───┤
                            │                       │                    │
                            │                       │<── Data loaded ────┤
                            │                       │                    │
                            │                       │ Initialize         │
                            │                       │ Services           │
                            │                       │                    │
                            │<── App Ready ─────────┤                    │
                            │                       │                    │
                     [ Application Running ]        │                    │
                            │                       │                    │
                            │ Shutdown              │ shutdown(ctx)      │
                            │──────────────────────>│                    │
                            │                       │──> store.Save() ───┤
                            │                       │<── Saved ──────────┤
                            │                       │                    │
```

### Backend Package Structure

```
btrack/
├── main.go                    # Wails entry point, app configuration
├── app.go                     # App struct, lifecycle, service initialization
│
├── Domain method handlers (exposed to frontend via Wails bindings):
├── project.go                 # Project CRUD operations
├── calculator.go              # Hour distribution algorithm
├── weekly_entry.go            # Weekly entry operations
├── dashboard.go               # Dashboard aggregations
├── meeting.go                 # Meeting operations
├── note.go                    # Note operations
├── goal.go                    # Goal operations
├── task.go                    # Task operations
├── template.go                # Template operations
├── backup.go                  # Backup/restore
├── export.go                  # CSV export
├── reports.go                 # Analytics/reports
│
├── internal/
│   ├── store/                 # In-memory state + JSON persistence (~1,780 lines)
│   │   ├── store.go           # Core: Load(), Save(), NextID()
│   │   ├── types.go           # Data, ProjectWithNested, Template
│   │   ├── projects.go        # Project CRUD (~260 lines)
│   │   ├── weekly_entries.go  # Weekly entry operations (~156 lines)
│   │   ├── nested_entities.go # Meetings, notes, goals, tasks (~535 lines)
│   │   ├── search.go          # Cross-project search (~250 lines)
│   │   ├── aggregations.go    # Dashboard, trends, reports (~336 lines)
│   │   └── templates.go       # Template CRUD (~74 lines)
│   │
│   ├── services/              # Business logic layer (~2,038 lines)
│   │   ├── project/
│   │   │   ├── service.go     # Project service, calculations
│   │   │   └── operations.go  # Update, delete, toggle operations
│   │   ├── tracking/
│   │   │   ├── service.go
│   │   │   ├── dashboard.go   # Dashboard data
│   │   │   ├── weekly.go      # Weekly entry operations
│   │   │   └── reports.go     # Monthly trends, variance, capacity
│   │   ├── notes/
│   │   │   ├── service.go
│   │   │   ├── meeting.go     # Meeting CRUD
│   │   │   ├── note.go        # Note CRUD
│   │   │   ├── goal.go        # Goal CRUD
│   │   │   └── task.go        # Task CRUD and filtering
│   │   └── system/
│   │       ├── service.go
│   │       ├── backup.go      # Backup/restore JSON file
│   │       ├── export.go      # CSV exports
│   │       └── template.go    # Template operations
│   │
│   └── models/                # Domain models (~893 lines)
│       ├── project.go         # Project, ProjectWithStats, ProjectHealth
│       ├── weekly_entry.go    # WeeklyEntry, WeeklyEntryWithStatus
│       ├── meeting.go         # Meeting, MeetingWithProject
│       ├── note.go            # Note, NoteWithProject
│       ├── goal.go            # Goal, GoalWithProject
│       ├── task.go            # Task, TaskWithContext
│       ├── apperror.go        # Structured error types
│       └── errors.go          # Error variables, NullableString helper
```

### App Struct

The `App` struct is the central point for all backend operations. All **exported** (capitalized) methods are automatically exposed to the frontend via Wails bindings.

```go
type App struct {
    ctx   context.Context
    store *store.Store
    
    // Services
    projectService  *project.Service
    trackingService *tracking.Service
    notesService    *notes.Service
    systemService   *system.Service
}

// Lifecycle methods
func (a *App) startup(ctx context.Context)  // Loads store, initializes services
func (a *App) shutdown(ctx context.Context) // Saves store

// Domain methods (examples - full list in WARP.md)
func (a *App) CreateProject(input CreateProjectInput) (*ProjectWithStats, error)
func (a *App) GetAllProjects(activeOnly bool) ([]ProjectWithStats, error)
func (a *App) GetDashboardData(weeksBack, weeksForward int) ([]DashboardWeekData, error)
func (a *App) UpdateActualHours(input UpdateActualHoursInput) (*WeeklyEntryWithStatus, error)
// ... ~40+ more methods
```

## Data Model

### JSON File Structure

All data is stored in a single JSON file with a project-centric nested hierarchy:

```json
{
  "schema_version": 1,
  "next_ids": {
    "project": 11,
    "weekly_entry": 100,
    "meeting": 26,
    "note": 5,
    "goal": 8,
    "task": 15,
    "template": 2
  },
  "projects": [
    {
      "id": 1,
      "name": "Client Project Alpha",
      "total_sold_hours": 160,
      "specialist_hours": 40,
      "start_date": "2026-01-01",
      "end_date": "2026-03-31",
      "is_active": true,
      "is_persistent": false,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-15T10:30:00Z",
      
      "weekly_entries": [
        {
          "id": 1,
          "project_id": 1,
          "week_start_date": "2025-12-30",
          "week_number": 1,
          "planned_hours": 12,
          "actual_hours": 10,
          "created_at": "2026-01-01T00:00:00Z",
          "updated_at": "2026-01-07T16:00:00Z"
        }
      ],
      
      "meetings": [
        {
          "id": 1,
          "project_id": 1,
          "title": "Kickoff Meeting",
          "meeting_date": "2026-01-05",
          "duration_minutes": 60,
          "attendees": "Alice, Bob, Charlie",
          "notes": "# Meeting Notes\n\n- Discussed timeline\n- Assigned tasks",
          "created_at": "2026-01-05T14:00:00Z",
          "updated_at": "2026-01-05T15:00:00Z"
        }
      ],
      
      "notes": [
        {
          "id": 1,
          "project_id": 1,
          "title": "Architecture Decisions",
          "content": "# Tech Stack\n\n- Frontend: React\n- Backend: Go",
          "created_at": "2026-01-10T10:00:00Z",
          "updated_at": "2026-01-10T11:00:00Z"
        }
      ],
      
      "goals": [
        {
          "id": 1,
          "project_id": 1,
          "title": "Complete Phase 1",
          "description": "Finish core features",
          "status": "in_progress",
          "target_date": "2026-02-15",
          "created_at": "2026-01-01T00:00:00Z",
          "updated_at": "2026-01-15T10:00:00Z"
        }
      ],
      
      "tasks": [
        {
          "id": 1,
          "project_id": 1,
          "source_type": "meeting",
          "source_id": 1,
          "title": "Set up repository",
          "description": "Initialize Git repo and CI/CD",
          "status": "completed",
          "priority": "high",
          "due_date": "2026-01-10",
          "created_at": "2026-01-05T15:00:00Z",
          "updated_at": "2026-01-08T09:00:00Z"
        }
      ]
    }
  ],
  
  "templates": [
    {
      "id": 1,
      "name": "Standard Engagement",
      "total_sold_hours": 160,
      "specialist_hours": 40,
      "created_at": "2026-01-01T00:00:00Z",
      "updated_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

### Entity Relationships

```
Projects (top-level array)
├── Weekly Entries (one-to-many, cascade delete)
├── Meetings (one-to-many, cascade delete)
├── Notes (one-to-many, cascade delete)
├── Goals (one-to-many, cascade delete)
└── Tasks (one-to-many, cascade delete)
    ├── Standalone (source_type = "standalone")
    ├── Linked to Meeting (source_type = "meeting", source_id = meeting.id)
    └── Linked to Note (source_type = "note", source_id = note.id)

Templates (top-level array, independent)
```

**Key Characteristics:**
- **Project-centric**: Each project contains all its related entities
- **Cascade deletion**: Deleting a project removes all nested data
- **No cross-project links**: All relationships stay within a project
- **Global IDs**: IDs are unique across all projects (via next_ids counters)

## Store Package

The `store` package manages all data operations with thread safety and atomic persistence.

### Core Operations

**Load**: Reads JSON file into memory
```go
store := store.New(filePath)
err := store.Load()  // Parses JSON into Data struct
```

**Save**: Writes state to disk atomically
```go
err := store.Save()  // Atomic: write to .tmp, then rename
```

**Thread Safety**: All operations use `sync.RWMutex`
- Read operations: `RLock()` / `RUnlock()` 
- Write operations: `Lock()` / `Unlock()` + automatic save

**NextID**: Auto-incrementing ID generation
```go
id := store.NextID("project")  // Returns next ID and increments counter
```

### Store Methods (examples)

```go
// Projects
CreateProject(input) (*Project, error)
GetProject(id) (*ProjectWithNested, error)
GetAllProjects(activeOnly) ([]ProjectWithNested, error)
UpdateProject(input) (*Project, error)
DeleteProject(id) error
SearchProjects(query) ([]Project, error)

// Weekly Entries
AddWeeklyEntry(projectID, entry) (*WeeklyEntry, error)
GetWeeklyEntries(projectID) ([]WeeklyEntry, error)
UpdateWeeklyEntry(projectID, entry) error

// Nested Entities (meetings, notes, goals, tasks)
AddMeeting(projectID, meeting) (*Meeting, error)
GetMeetings(projectID) ([]Meeting, error)
UpdateMeeting(projectID, meeting) error
DeleteMeeting(projectID, meetingID) error

// Aggregations
GetWeeklyAggregates(start, end) ([]WeeklyAggregate, error)
GetMonthlyTrends(monthsBack) ([]MonthlyTrend, error)
GetProjectStats(projectID) (totalPlanned, totalActual, totalWeeks, error)

// Templates
CreateTemplate(template) (*Template, error)
GetTemplates() ([]Template, error)
```

## Core Algorithms

### Frontloaded Hour Distribution

Distributes project hours across weeks with remainder frontloaded to earliest weeks:

```
Input: 
  TotalSoldHours = 160
  SpecialistHours = 40
  StartDate = 2026-01-01
  EndDate = 2026-03-31

Calculate:
  MyHours = 160 - 40 = 120
  Weeks = 13
  BaseHours = 120 / 13 = 9
  Remainder = 120 % 13 = 3

Distribute:
  Week 1: 9 + 1 = 10 hours
  Week 2: 9 + 1 = 10 hours  
  Week 3: 9 + 1 = 10 hours
  Week 4-13: 9 hours each
```

### Project Health Calculation

```go
// Persistent projects → always "on_track"
if project.IsPersistent {
    return "on_track"
}

// Completed projects
if now.After(endDate) {
    if actualHours <= plannedHours {
        return "completed"  // Success
    }
    return "over_budget"   // Completed but over
}

// Active projects
variance := actualHours / plannedHours
if variance > 1.2 {
    return "over_budget"   // >20% over
} else if variance > 1.0 {
    return "at_risk"       // Over planned
}
return "on_track"          // Normal
```

## Frontend Architecture

### Component Hierarchy

```
App.tsx (root)
└── AppProvider (context)
    └── AppContent
        └── Layout (navigation)
            ├── Dashboard
            │   ├── SummaryCards
            │   ├── WeeklyChart
            │   ├── WeeklyTable
            │   └── DailyAgenda
            │
            ├── WeekView
            │   ├── WeekHoursTable
            │   └── WeekMeetingList
            │
            ├── ProjectList
            │   └── ProjectCard (multiple)
            │
            ├── ProjectDetail
            │   ├── WeeklyBreakdown
            │   ├── GoalList
            │   ├── TaskList
            │   ├── NoteList
            │   └── MeetingList
            │
            ├── TasksView
            │   └── TaskList (filtered)
            │
            ├── ReportsView
            │   ├── TrendChart
            │   └── VarianceTable
            │
            └── SettingsView
                └── Backup/Export controls
```

### State Management Pattern

Uses the **useCrud factory hook** for consistent async operations:

```typescript
// Factory hook
function useCrud<T, CreateInput, UpdateInput>(config: {
  loadFn: () => Promise<T[]>,
  createFn: (input: CreateInput) => Promise<T>,
  updateFn: (input: UpdateInput) => Promise<T>,
  deleteFn: (id: number) => Promise<void>,
  getId: (item: T) => number
}) {
  const [items, setItems] = useState<T[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  
  // Methods: load(), create(), update(), remove()
  return { items, isLoading, error, load, create, update, remove }
}

// Usage in domain hooks
function useMeetings(projectId: number) {
  const config = useMemo(() => ({
    loadFn: () => GetMeetings(projectId),
    createFn: CreateMeeting,
    updateFn: UpdateMeeting,
    deleteFn: DeleteMeeting,
    getId: (m: Meeting) => m.id,
  }), [projectId])
  
  return useCrud<Meeting, CreateMeetingInput, UpdateMeetingInput>(config)
}
```

### Data Flow

```
Component
    ↓ calls
Custom Hook (e.g., useMeetings)
    ↓ uses
useCrud (factory)
    ↓ calls
Wails Bindings (frontend/wailsjs/go/main/App)
    ↓ Promise (JSON over IPC)
App Methods (Go)
    ↓ calls
Service Layer
    ↓ calls
Store (in-memory)
    ↓ persists to
JSON File
```

## Key Design Patterns

### 1. In-Memory Store with Atomic Persistence
All data lives in memory for fast access. Writes use temp-file-then-rename for atomicity:
```go
jsonData := marshal(s.data)
writeFile(filePath + ".tmp", jsonData)
rename(filePath + ".tmp", filePath)  // Atomic
```

### 2. Project-Centric Nesting
Projects contain all related entities. No foreign keys or joins needed:
```go
type ProjectWithNested struct {
    models.Project
    WeeklyEntries []models.WeeklyEntry
    Meetings      []models.Meeting
    Notes         []models.Note
    Goals         []models.Goal
    Tasks         []models.Task
}
```

### 3. Empty Slice Convention
Always return empty slices, never nil, to prevent frontend errors:
```go
if result == nil {
    result = []Meeting{}
}
return result, nil
```

### 4. Service Layer Pattern
Services encapsulate business logic and coordinate between store and domain methods:
```go
type Service struct {
    store *store.Store
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*Result, error) {
    // Validation
    // Store operations
    // Post-processing
    return result, nil
}
```

### 5. Structured Errors
Use `AppError` for consistent error handling:
```go
return nil, models.NotFound("project")
return nil, models.ValidationError("name", "name is required")
return nil, models.Internal(err, "failed to save")
```

## Security & Performance

### Security
- **Local-only storage**: No network exposure
- **Desktop app**: No server component
- **Native dialogs**: OS-level file pickers for backup/restore
- **Atomic writes**: Prevents data corruption during saves

### Performance
- **In-memory queries**: All reads from memory (no I/O)
- **Fast searches**: Linear scans acceptable for desktop workloads (<1000 projects)
- **Lazy loading**: Weekly entries loaded per-project
- **Efficient marshaling**: JSON with indentation for readability (107KB for 10 projects)
- **Thread-safe**: RWMutex allows concurrent reads

### Trade-offs
- **Memory usage**: All data in memory (negligible for typical usage)
- **Write amplification**: Entire file rewritten on changes (fast with atomic writes)
- **Single process**: No concurrent multi-user access (acceptable for desktop)
- **No schema migrations**: Just code changes to handle version upgrades

## Extension Points

**Add new entity type:**
1. Add model to `internal/models/`
2. Add CRUD methods to `internal/store/nested_entities.go`
3. Add service methods to appropriate service
4. Add domain methods to root-level `.go` file
5. Add frontend hook and components

**Add new report:**
1. Add aggregation method to `internal/store/aggregations.go`
2. Add service method to `internal/services/tracking/reports.go`
3. Add domain method to `reports.go`
4. Add frontend component in `components/reports/`

**Add new field to existing entity:**
1. Update model in `internal/models/`
2. Update JSON structure (auto-handled by Go's json tags)
3. Update validation if needed
4. Update frontend types and UI

No database migrations needed - the JSON structure evolves with the code.
