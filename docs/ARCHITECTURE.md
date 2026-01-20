# btrack Architecture Documentation

## Overview

btrack (Bandwidth Tracker) is a desktop application for tracking project bandwidth allocation across weeks. It follows a hybrid architecture pattern using **Wails v2**, which combines a Go backend with a React/TypeScript frontend, communicating through auto-generated bindings.

## High-Level Architecture

```mermaid
graph TB
    subgraph "Desktop Application"
        subgraph "Frontend (React/TypeScript)"
            UI[UI Components]
            Hooks[Custom Hooks]
            Context[App Context]
            Types[TypeScript Types]
        end
        
        subgraph "Wails Runtime"
            Bindings[Auto-generated Bindings]
            Events[Event System]
            Dialogs[Native Dialogs]
        end
        
        subgraph "Backend (Go)"
            App[App Struct]
            Domain[Domain Logic]
            Data[Data Layer]
        end
        
        subgraph "Storage"
            SQLite[(SQLite Database)]
        end
    end
    
    UI --> Hooks
    Hooks --> Bindings
    Context --> UI
    Bindings --> App
    App --> Domain
    Domain --> Data
    Data --> SQLite
    Dialogs --> App
```

## Technology Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| Frontend | React 18, TypeScript, Vite | UI rendering and state management |
| Styling | Tailwind CSS | Utility-first CSS framework |
| Charts | Recharts | Data visualization |
| Desktop Framework | Wails v2 | Desktop app packaging and Go-JS bridge |
| Backend | Go 1.24+ | Business logic and data operations |
| Database | SQLite (modernc.org/sqlite) | Local data persistence |

## Backend Architecture

### Application Lifecycle

```mermaid
sequenceDiagram
    participant Main as main.go
    participant Wails as Wails Runtime
    participant App as App Struct
    participant DB as Database
    
    Main->>Wails: wails.Run()
    Wails->>App: startup(ctx)
    App->>DB: database.Initialize()
    DB->>DB: RunMigrations()
    DB->>DB: SeedPersistentProjects()
    DB-->>App: *sql.DB
    Note over App: App ready for method calls
    
    Wails-->>Main: Application Running
    
    Note over Main,DB: ... Application Lifetime ...
    
    Wails->>App: shutdown(ctx)
    App->>DB: db.Close()
```

### Backend Package Structure

```mermaid
graph TB
    subgraph "Root Package (main)"
        main[main.go<br/>Entry point, Wails config]
        app[app.go<br/>App struct, lifecycle]
        
        subgraph "Domain Files"
            project[project.go<br/>Project CRUD]
            calculator[calculator.go<br/>Frontloading algorithm]
            weekly[weekly_entry.go<br/>Weekly entries]
            dashboard[dashboard.go<br/>Aggregations]
            meeting[meeting.go<br/>Meeting CRUD]
            note[note.go<br/>Note CRUD]
            goal[goal.go<br/>Goal CRUD]
            task[task.go<br/>Task CRUD]
            template[template.go<br/>Templates]
            backup[backup.go<br/>Backup/Restore]
            export[export.go<br/>CSV Export]
            reports[reports.go<br/>Analytics]
        end
    end
    
    subgraph "internal/database"
        database[database.go<br/>Connection setup]
        migrations[migrations.go<br/>Schema definitions]
        queries[queries.go<br/>SQL constants]
        scanner[scanner.go<br/>Helper utilities]
    end
    
    subgraph "internal/models"
        projectModel[project.go<br/>Project types]
        weeklyModel[weekly_entry.go<br/>Entry types]
        meetingModel[meeting.go<br/>Meeting types]
        noteModel[note.go<br/>Note types]
        goalModel[goal.go<br/>Goal types]
        taskModel[task.go<br/>Task types]
        errors[errors.go<br/>Error types]
    end
    
    app --> database
    project --> queries
    project --> projectModel
    meeting --> meetingModel
    meeting --> scanner
```

### App Struct & Method Binding

The `App` struct is the central point for all backend operations. All **exported** (capitalized) methods are automatically exposed to the frontend via Wails bindings.

```mermaid
classDiagram
    class App {
        -ctx context.Context
        -db *sql.DB
        +startup(ctx)
        +shutdown(ctx)
        +CreateProject(input) ProjectWithStats
        +GetAllProjects(activeOnly) []ProjectWithStats
        +GetProject(id) ProjectWithStats
        +UpdateProject(input) ProjectWithStats
        +DeleteProject(id) error
        +CalculateDistribution(input) []WeeklyEntry
        +GetWeeklyEntries(projectID) []WeeklyEntryWithStatus
        +UpdateActualHours(input) WeeklyEntryWithStatus
        +GetDashboardData(weeksBack, weeksForward) []DashboardWeekData
        +GetDashboardSummary() DashboardSummary
        +CreateMeeting(input) Meeting
        +GetMeetings(projectID) []Meeting
        +CreateNote(input) Note
        +GetNotes(projectID) []Note
        +CreateGoal(input) Goal
        +GetGoals(projectID) []Goal
        +CreateTask(input) Task
        +GetAllTasks(status, projectID) []TaskWithContext
        +CreateBackup() string
        +RestoreBackup() error
        +ExportWeeklyReport(start, end) string
        +GetMonthlyTrends(monthsBack) []MonthlyTrend
    }
```

## Data Model

### Entity Relationship Diagram

```mermaid
erDiagram
    PROJECTS ||--o{ WEEKLY_ENTRIES : contains
    PROJECTS ||--o{ PROJECT_MEETINGS : has
    PROJECTS ||--o{ PROJECT_NOTES : has
    PROJECTS ||--o{ PROJECT_GOALS : has
    PROJECTS ||--o{ TASKS : has
    PROJECT_MEETINGS ||--o{ TASKS : spawns
    PROJECT_NOTES ||--o{ TASKS : spawns
    
    PROJECTS {
        int id PK
        string name
        int total_sold_hours
        int specialist_hours
        string start_date
        string end_date
        bool is_active
        bool is_persistent
        datetime created_at
        datetime updated_at
    }
    
    WEEKLY_ENTRIES {
        int id PK
        int project_id FK
        string week_start_date
        int week_number
        int planned_hours
        int actual_hours
        datetime created_at
        datetime updated_at
    }
    
    PROJECT_MEETINGS {
        int id PK
        int project_id FK
        string title
        string meeting_date
        int duration_minutes
        string attendees
        string notes
        datetime created_at
        datetime updated_at
    }
    
    PROJECT_NOTES {
        int id PK
        int project_id FK
        string title
        string content
        datetime created_at
        datetime updated_at
    }
    
    PROJECT_GOALS {
        int id PK
        int project_id FK
        string title
        string description
        string status
        string target_date
        datetime created_at
        datetime updated_at
    }
    
    TASKS {
        int id PK
        int project_id FK
        string source_type
        int source_id
        string title
        string description
        string status
        string priority
        string due_date
        datetime created_at
        datetime updated_at
    }
    
    PROJECT_TEMPLATES {
        int id PK
        string name
        int total_sold_hours
        int specialist_hours
        datetime created_at
        datetime updated_at
    }
```

### Key Data Relationships

1. **Projects → Weekly Entries**: One-to-many with cascade delete
2. **Projects → Meetings/Notes/Goals**: One-to-many with cascade delete
3. **Tasks → Sources**: Polymorphic relationship via `source_type` ('standalone', 'meeting', 'note')
4. **Templates**: Standalone, no foreign keys

## Core Algorithms

### Frontloaded Hour Distribution

The `CalculateDistribution` function implements a frontloading algorithm to distribute project hours across weeks:

```mermaid
flowchart TD
    A[Input: TotalSoldHours, SpecialistHours, StartDate, EndDate] --> B[Calculate MyHours = TotalSoldHours - SpecialistHours]
    B --> C[Calculate Weeks between dates]
    C --> D[BaseHours = MyHours / Weeks]
    D --> E[Remainder = MyHours % Weeks]
    E --> F{For each week}
    F --> G{Week index < Remainder?}
    G -->|Yes| H[Assign BaseHours + 1]
    G -->|No| I[Assign BaseHours]
    H --> J[Next week]
    I --> J
    J --> F
    F -->|Done| K[Return Weekly Entries]
```

**Example**: 10 MyHours / 3 Weeks → Week 1: 4 hours, Week 2: 3 hours, Week 3: 3 hours

### Project Health Calculation

```mermaid
flowchart TD
    A[Project with Stats] --> B{Is Persistent?}
    B -->|Yes| C[Return: on_track]
    B -->|No| D{Past End Date?}
    D -->|Yes| E[Return: completed]
    D -->|No| F{Actual > Planned?}
    F -->|Yes| G[Return: over_budget]
    F -->|No| H{Actual > 80% Planned?}
    H -->|Yes| I[Return: at_risk]
    H -->|No| J[Return: on_track]
```

## Frontend Architecture

### Component Hierarchy

```mermaid
graph TB
    App[App.tsx] --> AppProvider[AppProvider]
    AppProvider --> AppContent[AppContent]
    AppContent --> Layout[Layout]
    
    Layout --> Dashboard[Dashboard]
    Layout --> WeekView[WeekView]
    Layout --> ProjectList[ProjectList]
    Layout --> ProjectDetail[ProjectDetail]
    Layout --> TasksView[TasksView]
    Layout --> ReportsView[ReportsView]
    Layout --> SettingsView[SettingsView]
    
    Dashboard --> DashboardCharts[Charts]
    Dashboard --> DashboardStats[Summary Stats]
    
    ProjectDetail --> ProjectMeetings[Meetings]
    ProjectDetail --> ProjectNotes[Notes]
    ProjectDetail --> ProjectGoals[Goals]
    ProjectDetail --> WeeklyEntries[Weekly Entries]
    
    subgraph "Shared UI Components"
        EntityList[EntityList]
        Modal[Modal]
        Card[Card]
        Button[Button]
    end
```

### State Management Pattern

```mermaid
flowchart LR
    subgraph "React Component"
        Component[Component]
    end
    
    subgraph "Custom Hook"
        Hook[useProjects / useMeetings / etc.]
        CrudHook[useCrud Factory]
    end
    
    subgraph "Wails Bindings"
        Bindings[frontend/wailsjs/go/main/App]
    end
    
    subgraph "Go Backend"
        Backend[App Methods]
    end
    
    Component --> Hook
    Hook --> CrudHook
    CrudHook --> Bindings
    Bindings -->|Promise| Backend
    Backend -->|JSON| Bindings
```

### Custom Hooks Architecture

The `useCrud` factory hook provides standardized async state management:

```mermaid
classDiagram
    class useCrud~T,CreateInput,UpdateInput~ {
        +items: T[]
        +isLoading: boolean
        +error: string
        +load(): Promise
        +create(input: CreateInput): Promise
        +update(input: UpdateInput): Promise
        +remove(id: number): Promise
    }
    
    class useProjects {
        +projects: ProjectWithStats[]
        +activeProjects: ProjectWithStats[]
        +createProject()
        +updateProject()
        +deleteProject()
    }
    
    class useMeetings {
        +meetings: Meeting[]
        +createMeeting()
        +updateMeeting()
        +deleteMeeting()
    }
    
    useProjects --> useCrud
    useMeetings --> useCrud
```

## Data Flow

### Create Project Flow

```mermaid
sequenceDiagram
    participant UI as ProjectForm
    participant Hook as useProjects
    participant Wails as Wails Bindings
    participant Go as App.CreateProject
    participant DB as SQLite
    
    UI->>Hook: handleSubmit(formData)
    Hook->>Wails: CreateProject(input)
    Wails->>Go: CreateProject(models.CreateProjectInput)
    Go->>Go: input.Validate()
    Go->>Go: CalculateDistribution(input)
    Go->>DB: BEGIN TRANSACTION
    Go->>DB: INSERT INTO projects
    Go->>DB: INSERT INTO weekly_entries (for each week)
    Go->>DB: COMMIT
    Go->>Go: GetProject(projectID)
    Go-->>Wails: ProjectWithStats
    Wails-->>Hook: Promise resolves
    Hook->>Hook: Update state
    Hook-->>UI: Re-render with new project
```

### Weekly View Data Flow

```mermaid
sequenceDiagram
    participant UI as WeekView
    participant Hook as useWeekView
    participant Wails as Wails Bindings
    participant Go as Backend Methods
    participant DB as SQLite
    
    UI->>Hook: useWeekView()
    Hook->>Wails: GetWeeklyEntriesByWeek(weekStartDate)
    Wails->>Go: GetWeeklyEntriesByWeek(weekStartDate)
    Go->>Go: ensurePersistentEntries()
    Go->>DB: INSERT OR IGNORE persistent entries
    Go->>DB: SELECT entries with project names
    Go-->>Wails: []WeeklyEntryWithProject
    Wails-->>Hook: Promise resolves
    
    Hook->>Wails: GetMeetingsByWeek(weekStartDate)
    Wails->>Go: GetMeetingsByWeek(weekStartDate)
    Go->>DB: SELECT meetings for date range
    Go-->>Wails: []MeetingWithProject
    Wails-->>Hook: Promise resolves
    
    Hook-->>UI: { entries, meetings, isLoading }
```

## File Organization

```
btrack/
├── main.go                    # Wails entry point
├── app.go                     # App struct, lifecycle
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
│   ├── database/
│   │   ├── database.go        # Connection management
│   │   ├── migrations.go      # Schema definitions
│   │   ├── queries.go         # SQL query constants
│   │   └── scanner.go         # Helper utilities
│   │
│   └── models/
│       ├── project.go         # Project domain models
│       ├── weekly_entry.go    # Entry models
│       ├── meeting.go         # Meeting models
│       ├── note.go            # Note models
│       ├── goal.go            # Goal models
│       ├── task.go            # Task models
│       ├── search.go          # Search result models
│       └── errors.go          # Custom error types
│
├── frontend/
│   ├── src/
│   │   ├── App.tsx            # Root component
│   │   ├── main.tsx           # React entry point
│   │   ├── components/
│   │   │   ├── dashboard/     # Dashboard views
│   │   │   ├── projects/      # Project views
│   │   │   ├── week/          # Week view
│   │   │   ├── weekly/        # Weekly entry components
│   │   │   ├── meetings/      # Meeting components
│   │   │   ├── notes/         # Note components
│   │   │   ├── goals/         # Goal components
│   │   │   ├── tasks/         # Task components
│   │   │   ├── reports/       # Report views
│   │   │   ├── settings/      # Settings views
│   │   │   ├── layout/        # Layout components
│   │   │   └── ui/            # Shared UI components
│   │   ├── hooks/             # Custom React hooks
│   │   ├── context/           # React context providers
│   │   ├── types/             # TypeScript definitions
│   │   └── utils/             # Utility functions
│   │
│   └── wailsjs/               # Auto-generated (DO NOT EDIT)
│       └── go/main/App.ts     # Go method bindings
│
├── build/                     # Build configuration
├── wails.json                 # Wails project config
├── go.mod                     # Go dependencies
└── frontend/package.json      # Node dependencies
```

## Key Design Patterns

### 1. Repository Pattern (implicit)
SQL queries are centralized in `internal/database/queries.go`, providing a single source for all database operations.

### 2. Domain Model Pattern
Business entities are defined in `internal/models/` with validation methods and computed properties.

### 3. Scanner Functions
Entity models include `Scan` functions that encapsulate row scanning logic:
```go
// ScanMeeting scans a database row into a Meeting struct
meeting, err := models.ScanMeeting(rows.Scan)
```

### 4. Slice Helpers
Prevent nil JSON responses using `database.EnsureSlice()`:
```go
return database.EnsureSlice(meetings), nil
```

### 5. Factory Hook Pattern (Frontend)
The `useCrud` hook provides reusable CRUD state management that domain-specific hooks build upon.

## Security Considerations

1. **Local-only storage**: All data is stored locally in SQLite
2. **No network exposure**: Desktop app with no server component
3. **Native dialogs**: Backup/restore uses OS-native file dialogs
4. **Cascading deletes**: Foreign key constraints with CASCADE ensure referential integrity

## Performance Considerations

1. **Indexed queries**: Key columns are indexed (project_id, week_start_date, status)
2. **Lazy loading**: Weekly entries loaded per-project, not all at once
3. **Pagination**: Search results limited to 50 items
4. **Embedded assets**: Frontend built and embedded into binary via `//go:embed`

## Extension Points

1. **New entities**: Add model, migration, queries, domain file, and frontend hook
2. **Reports**: Add query in `queries.go`, method in `reports.go`
3. **Export formats**: Add methods in `export.go` following existing CSV pattern
4. **Frontend views**: Add component folder, hook, and route in `App.tsx`
