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
- **Templates**: Save projects as reusable templates for quick project creation with predefined hour allocations

### Data Management
- **Search**: Full-text search across projects, meetings, notes, and goals
- **Backup & Restore**: Create backups and restore from previous database states with safety mechanisms
- **Export**: Export data to CSV format (weekly reports, project summaries, or all projects)

### Analytics & Reporting
- **Monthly Trends**: View planned vs actual hours with variance analysis over multiple months
- **Variance Reports**: Analyze planned vs actual hours by project for any date range
- **Capacity Utilization**: Track weekly capacity utilization with percentage calculations
- **Project Health**: Automatic health status tracking (on track/at risk/over budget/completed)

## Architecture

- **Backend**: Go with SQLite (modernc.org/sqlite driver)
- **Frontend**: React 18 + TypeScript + Vite + Tailwind CSS + Recharts
- **Framework**: Wails v2 (native desktop apps with web technologies)
- **Database**: Local SQLite at `~/Library/Application Support/btrack/btrack.db`

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
├── template.go             # Project template management
├── backup.go               # Database backup and restore
├── export.go               # CSV export functionality
├── reports.go              # Analytics and reporting
├── internal/
│   ├── database/           # SQLite setup, migrations, queries
│   └── models/             # Domain models and validation
└── frontend/
    ├── src/
    │   ├── components/     # React components
    │   ├── hooks/          # Custom React hooks
    │   ├── context/        # React context providers
    │   ├── types/          # TypeScript definitions
    │   └── utils/          # Utility functions
    └── wailsjs/            # Auto-generated Wails bindings (DO NOT EDIT)
```

## Database

### Location
The SQLite database is stored at:
```
~/Library/Application Support/btrack/btrack.db
```

### Reset Database (Development)
To reset the database schema during development:

```bash
rm ~/Library/Application\ Support/btrack/btrack.db
```

The database will be recreated with fresh schema and seeded persistent projects ("Management" and "Internal Projects") on next run.

### Backup & Restore
Use the built-in backup/restore features in the application:
- **Backup**: Creates a timestamped copy of the database to a location you choose
- **Restore**: Restores from a backup file with automatic safety backup (`.pre-restore`) in case of failure

## More Information

- See [WARP.md](WARP.md) for detailed architecture and development guidance
- See [CLAUDE.md](CLAUDE.md) for AI assistant development guidelines
- Wails documentation: https://wails.io/docs
