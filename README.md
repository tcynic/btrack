# btrack - Bandwidth Tracker

A desktop application for tracking project bandwidth allocation across weeks. Built with Wails v2, combining a Go backend with a React/TypeScript frontend.

## Features

### Core Functionality
- **Project Management**: Track multiple client projects with total sold hours and specialist hours allocation
- **Frontloaded Distribution**: Automatically distributes "My Hours" (Total - Specialist) across project timelines with frontloading to earliest weeks
- **Weekly Tracking**: Plan and track actual hours worked on a per-week basis (Monday start dates)
- **Dashboard Analytics**: View aggregated statistics across all active projects with health status indicators
- **Persistent Projects**: Special "Management" and "Internal Projects" for ongoing work allocation

### Project Features
- **Meetings**: Schedule and track project meetings with attendees, duration, and notes
- **Notes**: Create markdown-formatted notes for each project
- **Goals**: Set and track project goals with status management (pending/in progress/completed/cancelled) and completion rate tracking
- **Tasks**: Create and manage action items linked to meetings, notes, or as standalone tasks with priority and status tracking
- **Templates**: Save projects as reusable templates for quick project creation with predefined hour allocations

### Task Management
- **Aggregated Task View**: View all tasks across projects in a single, filterable list
- **Linked Tasks**: Create tasks from meetings or notes to track action items in context
- **Status Tracking**: Quick status updates (pending → in progress → completed) via clickable badges
- **Priority Levels**: Organize tasks by low, medium, or high priority
- **Filtering**: Filter tasks by status, project, or due date to focus on what matters

### Data Management
- **Search**: Full-text search across projects, meetings, notes, goals, and tasks
- **Backup & Restore**: Create backups and restore from previous database states with safety mechanisms
- **Export**: Export data to CSV format (weekly reports, project summaries, or all projects)

### Analytics & Reporting
- **Monthly Trends**: View planned vs actual hours with variance analysis over multiple months
- **Variance Reports**: Analyze planned vs actual hours by project for any date range
- **Capacity Utilization**: Track weekly capacity utilization with percentage calculations
- **Project Health**: Automatic health status tracking (on track/at risk/over budget/completed)

## Architecture

- **Backend**: Go with JSON storage
- **Frontend**: React 18 + TypeScript + Vite + Tailwind CSS + Recharts
- **Framework**: Wails v2 (native desktop apps with web technologies)
- **Data Storage**: Local JSON at `~/Library/Application Support/btrack/btrack-data.json`

## Development

### Prerequisites

- Go 1.24+
- Node.js/npm
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Live Development

```bash
wails dev
```

This runs both the Go backend and Vite dev server with hot reload. The frontend is accessible at the dev server URL, and Go methods are available at http://localhost:34115 for browser-based development.

### Building

```bash
wails build
```

Creates a production-ready native application binary for your platform.

### Frontend Commands

```bash
cd frontend
npm install          # Install dependencies
npm run dev          # Run Vite dev server independently
npm run build        # Build frontend assets
```

## Project Structure

```
.
├── main.go                 # Application entry point
├── app.go                  # Main App struct with lifecycle methods
├── project.go              # Project CRUD and search operations
├── calculator.go           # Frontloading distribution algorithm
├── weekly_entry.go         # Weekly entry operations
├── dashboard.go            # Dashboard aggregations
├── meeting.go              # Meeting management and search
├── note.go                 # Note management and search
├── goal.go                 # Goal management, statistics, and search
├── task.go                 # Task management, filtering, and search
├── template.go             # Project template management
├── backup.go               # Database backup and restore
├── export.go               # CSV export functionality
├── reports.go              # Analytics and reporting
├── internal/
│   ├── store/
│   │   ├── store.go           # Store with Load/Save and atomic writes
│   │   ├── types.go           # Data structures (ProjectWithNested, etc.)
│   │   ├── projects.go        # Project CRUD operations
│   │   ├── weekly_entries.go  # Weekly entry operations
│   │   ├── nested_entities.go # Meetings, notes, goals, tasks CRUD
│   │   ├── search.go          # Cross-project search
│   │   ├── aggregations.go    # Dashboard, reports, statistics
│   │   └── templates.go       # Template operations
│   ├── services/
│   │   ├── project/           # Project service layer
│   │   ├── tracking/          # Tracking service (dashboard, reports)
│   │   ├── notes/             # Notes service (meetings, notes, goals, tasks)
│   │   └── system/            # System service (backup, export, templates)
│   └── models/             # Domain models with validation
└── frontend/
    ├── src/
    │   ├── components/
    │   │   ├── ui/         # Reusable UI components (Button, Modal, EntityList, etc.)
    │   │   ├── meetings/   # Meeting-specific components
    │   │   ├── notes/      # Note-specific components
    │   │   ├── goals/      # Goal-specific components
    │   │   ├── tasks/      # Task-specific components
    │   │   └── ...         # Other feature components
    │   ├── hooks/
    │   │   ├── useCrud.ts  # Generic CRUD hook factory
    │   │   ├── useMeetings.ts, useNotes.ts, useGoals.ts  # Entity-specific hooks
    │   │   └── ...         # Other custom hooks
    │   ├── context/        # React context providers
    │   ├── types/          # TypeScript definitions
    │   └── utils/          # Utility functions
    └── wailsjs/            # Auto-generated Wails bindings (DO NOT EDIT)
```

## Code Patterns

### Backend (Go)

**Store Operations**: All CRUD operations go through the Store, which handles locking and persistence:
```go
// Store automatically locks, mutates, saves, and unlocks
project, err := s.store.CreateProject(input)
if err != nil {
  return nil, err
}
```

**Thread Safety**: Store uses sync.RWMutex for concurrent read/write safety.

**Empty Slice Handling**: Methods return empty slices instead of nil to prevent frontend errors:
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

## Data Storage

### Location
The JSON data file is stored at:
```
~/Library/Application Support/btrack/btrack-data.json
```

### Reset Data (Development)
To reset the data during development:

```bash
rm ~/Library/Application\ Support/btrack/btrack-data.json
```

The file will be recreated with empty data on next run. Persistent projects can be seeded via the store.

### Backup & Restore
Use the built-in backup/restore features in the application:
- **Backup**: Creates a timestamped copy of the JSON file to a location you choose
- **Restore**: Restores from a backup JSON file with automatic safety backup (`.pre-restore`) in case of failure

## More Information

- See [WARP.md](WARP.md) for detailed architecture and development guidance
- See [CLAUDE.md](CLAUDE.md) for AI assistant development guidelines
- Wails documentation: https://wails.io/docs
