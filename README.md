# btrack - Bandwidth Tracker

A desktop application for tracking project bandwidth allocation across weeks. Built with Wails v2, combining a Go backend with a React/TypeScript frontend.

## Features

- **Project Management**: Track multiple client projects with total sold hours and specialist hours allocation
- **Frontloaded Distribution**: Automatically distributes "My Hours" (Total - Specialist) across project timelines with frontloading to earliest weeks
- **Weekly Tracking**: Plan and track actual hours worked on a per-week basis (Monday start dates)
- **Dashboard Analytics**: View aggregated statistics across all active projects
- **Meetings**: Schedule and track project meetings with attendees and notes
- **Notes**: Create markdown-formatted notes for each project
- **Goals**: Set and track project goals with status management (pending/in progress/completed/cancelled)
- **Persistent Projects**: Special "Management" and "Internal Projects" for ongoing work allocation

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
├── project.go              # Project CRUD operations
├── calculator.go           # Frontloading distribution algorithm
├── weekly_entry.go         # Weekly entry operations
├── dashboard.go            # Dashboard aggregations
├── meeting.go              # Meeting management
├── note.go                 # Note management
├── goal.go                 # Goal management
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

## Database Reset

For development, to reset the database schema:

```bash
rm ~/Library/Application\ Support/btrack/btrack.db
```

The database will be recreated with fresh schema on next run.

## More Information

- See [WARP.md](WARP.md) for detailed architecture and development guidance
- See [CLAUDE.md](CLAUDE.md) for AI assistant development guidelines
- Wails documentation: https://wails.io/docs
