package repository

import (
	"context"
	"database/sql"

	"btrack/internal/models"
)

// Repository provides base database access functionality
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// WithTx executes a function within a database transaction
// Automatically handles commit/rollback based on error returns
func (r *Repository) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return models.DatabaseError(err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return models.DatabaseError(err)
	}

	return nil
}

// QuerySlice executes a query and scans results into a slice using a scan function
// Automatically handles EnsureSlice to prevent nil JSON responses
func (r *Repository) QuerySlice[T any](
	query string,
	args []any,
	scanFn func(func(dest ...any) error) (*T, error),
) ([]T, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, models.DatabaseError(err)
	}
	defer rows.Close()

	var items []T
	for rows.Next() {
		item, err := scanFn(rows.Scan)
		if err != nil {
			return nil, models.DatabaseError(err)
		}
		items = append(items, *item)
	}

	// Ensure we return empty slice instead of nil
	if items == nil {
		items = []T{}
	}

	return items, nil
}

// QueryOne executes a query and scans a single result
func (r *Repository) QueryOne[T any](
	query string,
	args []any,
	scanFn func(func(dest ...any) error) (*T, error),
) (*T, error) {
	row := r.db.QueryRow(query, args...)
	item, err := scanFn(row.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.NotFound("record")
		}
		return nil, models.DatabaseError(err)
	}
	return item, nil
}

// Exec executes a query that doesn't return rows
func (r *Repository) Exec(query string, args ...any) (sql.Result, error) {
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, models.DatabaseError(err)
	}
	return result, nil
}

// DB returns the underlying database connection
func (r *Repository) DB() *sql.DB {
	return r.db
}
