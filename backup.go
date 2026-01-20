package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
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

	// Auto-backup before restore
	if err := a.autoBackup(); err != nil {
		// Log warning but don't fail the operation
		log.Printf("Warning: auto-backup failed before restore: %v", err)
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

// autoBackup creates an automatic backup before destructive operations
func (a *App) autoBackup() error {
	dbPath := a.getDBPath()
	backupDir := filepath.Join(filepath.Dir(dbPath), "backups")
	
	// Create backups directory if it doesn't exist
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	
	// Create backup with timestamp
	backupName := fmt.Sprintf("auto-backup-%s.db", time.Now().Format("2006-01-02-150405"))
	backupPath := filepath.Join(backupDir, backupName)
	
	if err := copyFile(dbPath, backupPath); err != nil {
		return fmt.Errorf("failed to create auto-backup: %w", err)
	}
	
	log.Printf("Auto-backup created: %s", backupPath)
	
	// Clean up old backups
	if err := a.cleanOldBackups(backupDir, 5); err != nil {
		log.Printf("Warning: failed to clean old backups: %v", err)
	}
	
	return nil
}

// cleanOldBackups removes old auto-backups, keeping only the most recent N backups
func (a *App) cleanOldBackups(backupDir string, keepCount int) error {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}
	
	// Filter auto-backups and get their info
	type backupFile struct {
		name    string
		modTime time.Time
	}
	
	var backups []backupFile
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".db" {
			continue
		}
		
		// Only process auto-backups
		if len(entry.Name()) < 12 || entry.Name()[:12] != "auto-backup-" {
			continue
		}
		
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		backups = append(backups, backupFile{
			name:    entry.Name(),
			modTime: info.ModTime(),
		})
	}
	
	// Sort by modification time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].modTime.After(backups[j].modTime)
	})
	
	// Delete old backups beyond keepCount
	for i := keepCount; i < len(backups); i++ {
		backupPath := filepath.Join(backupDir, backups[i].name)
		if err := os.Remove(backupPath); err != nil {
			log.Printf("Warning: failed to delete old backup %s: %v", backups[i].name, err)
		} else {
			log.Printf("Deleted old backup: %s", backups[i].name)
		}
	}
	
	return nil
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
