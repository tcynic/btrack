package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"btrack/internal/database"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// BackupInfo provides information about the current database
type BackupInfo struct {
	DatabasePath string `json:"databasePath"`
	DatabaseSize int64  `json:"databaseSize"` // in bytes
	LastModified string `json:"lastModified"`
}

// GetBackupInfo returns information about the current database
func (a *App) GetBackupInfo() (*BackupInfo, error) {
	dbPath := a.getDBPath()
	
	info, err := os.Stat(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get database info: %w", err)
	}

	return &BackupInfo{
		DatabasePath: dbPath,
		DatabaseSize: info.Size(),
		LastModified: info.ModTime().Format("2006-01-02 15:04:05"),
	}, nil
}

// CreateBackup creates a backup of the database to a user-selected location
func (a *App) CreateBackup() (string, error) {
	ctx := context.Background()
	
	// Prompt user to select backup location
	defaultName := fmt.Sprintf("btrack-backup-%s.db", time.Now().Format("2006-01-02-150405"))
	savePath, err := runtime.SaveFileDialog(ctx, runtime.SaveDialogOptions{
		DefaultFilename: defaultName,
		Title:           "Save Database Backup",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Database Files (*.db)",
				Pattern:     "*.db",
			},
		},
	})
	
	if err != nil {
		return "", fmt.Errorf("failed to show save dialog: %w", err)
	}
	
	// User cancelled
	if savePath == "" {
		return "", nil
	}

	// Get source database path
	dbPath := a.getDBPath()
	
	// Copy database file
	if err := copyFile(dbPath, savePath); err != nil {
		return "", fmt.Errorf("failed to copy database: %w", err)
	}

	return savePath, nil
}

// RestoreBackup restores the database from a backup file
func (a *App) RestoreBackup() error {
	ctx := context.Background()
	
	// Prompt user to select backup file
	backupPath, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select Backup File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Database Files (*.db)",
				Pattern:     "*.db",
			},
		},
	})
	
	if err != nil {
		return fmt.Errorf("failed to show open dialog: %w", err)
	}
	
	// User cancelled
	if backupPath == "" {
		return nil
	}

	// Validate backup file exists
	if _, err := os.Stat(backupPath); err != nil {
		return fmt.Errorf("backup file not found: %w", err)
	}

	// Get current database path
	dbPath := a.getDBPath()
	
	// Create a temporary backup of current database before restore
	tempBackup := dbPath + ".pre-restore"
	if err := copyFile(dbPath, tempBackup); err != nil {
		return fmt.Errorf("failed to create temporary backup: %w", err)
	}

	// Close current database connection
	if err := a.db.Close(); err != nil {
		// Restore from temp backup if close fails
		os.Remove(tempBackup)
		return fmt.Errorf("failed to close database: %w", err)
	}

	// Copy backup file to database location
	if err := copyFile(backupPath, dbPath); err != nil {
		// Try to restore original database
		copyFile(tempBackup, dbPath)
		os.Remove(tempBackup)
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// Remove temporary backup
	os.Remove(tempBackup)

	// Reopen database connection
	db, err := database.Initialize()
	if err != nil {
		return fmt.Errorf("failed to reopen database: %w", err)
	}
	a.db = db

	return nil
}

// getDBPath returns the path to the database file
func (a *App) getDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, "Library", "Application Support", "btrack", "btrack.db")
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy contents
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Sync to ensure data is written
	return destFile.Sync()
}
