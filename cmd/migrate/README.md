# SQLite to JSON Migration Tool

This tool migrates data from the legacy SQLite database to the new JSON-based storage format.

## Usage

```bash
# Build the tool
go build -o migrate-tool ./cmd/migrate

# Preview migration (dry run)
./migrate-tool --dry-run

# Execute migration
./migrate-tool

# Custom paths
./migrate-tool --db /path/to/btrack.db --json /path/to/output.json
```

## Flags

- `--db <path>` - Source SQLite database path (default: auto-detected for macOS)
- `--json <path>` - Target JSON file path (default: auto-detected for macOS)
- `--dry-run` - Preview migration without writing files

## What it migrates

- Projects (with all fields)
- Weekly entries (linked to projects)
- Meetings (linked to projects)
- Notes (linked to projects)
- Goals (linked to projects)
- Tasks (linked to projects or standalone)
- Templates

## Safety

- Automatically creates timestamped backup before overwriting existing JSON
- Validates all data during migration
- Calculates next IDs from maximum existing IDs

## After Migration

Once migration is complete:
1. Verify data loaded correctly by running the application
2. Test all CRUD operations
3. Keep the SQLite database as a backup
4. The application will now use `btrack-data.json` for all operations
