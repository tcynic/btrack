package system

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// BackupInfo provides information about the current data file
type BackupInfo struct {
	DatabasePath string `json:"databasePath"`
	DatabaseSize int64  `json:"databaseSize"` // in bytes
	LastModified string `json:"lastModified"`
}

// GetBackupInfo returns information about the current data file
func (s *Service) GetBackupInfo() (*BackupInfo, error) {
	dbPath := getDBPath()
	
	info, err := os.Stat(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get data file info: %w", err)
	}

	return &BackupInfo{
		DatabasePath: dbPath,
		DatabaseSize: info.Size(),
		LastModified: info.ModTime().Format("2006-01-02 15:04:05"),
	}, nil
}

// CreateBackup creates a backup of the data file to a user-selected location
func (s *Service) CreateBackup() (string, error) {
	// Prompt user to select backup location
	defaultName := fmt.Sprintf("btrack-backup-%s.json", time.Now().Format("2006-01-02-150405"))
	savePath, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		DefaultFilename: defaultName,
		Title:           "Save Data Backup",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "JSON Files (*.json)",
				Pattern:     "*.json",
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

	// Get source data file path
	dbPath := getDBPath()
	
	// Copy data file
	if err := copyFile(dbPath, savePath); err != nil {
		return "", fmt.Errorf("failed to copy data file: %w", err)
	}

	return savePath, nil
}

// RestoreBackup restores the data file from a backup (store is automatically reloaded)
func (s *Service) RestoreBackup() error {
	// Prompt user to select backup file
	backupPath, err := runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Title: "Select Backup File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "JSON Files (*.json)",
				Pattern:     "*.json",
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
	if err := s.AutoBackup(); err != nil {
		log.Printf("Warning: auto-backup failed before restore: %v", err)
	}

	// Get current data file path
	dbPath := getDBPath()
	
	// Create a temporary backup of current data file before restore
	tempBackup := dbPath + ".pre-restore"
	if err := copyFile(dbPath, tempBackup); err != nil {
		return fmt.Errorf("failed to create temporary backup: %w", err)
	}

	// Copy backup file to data file location
	if err := copyFile(backupPath, dbPath); err != nil {
		// Try to restore original data file
		copyFile(tempBackup, dbPath)
		os.Remove(tempBackup)
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	// Remove temporary backup
	os.Remove(tempBackup)

	// Reload store from restored file
	if err := s.store.Load(); err != nil {
		return fmt.Errorf("failed to reload data: %w", err)
	}

	return nil
}

// AutoBackup creates an automatic backup before destructive operations
func (s *Service) AutoBackup() error {
	dbPath := getDBPath()
	backupDir := filepath.Join(filepath.Dir(dbPath), "backups")
	
	// Create backups directory if it doesn't exist
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	
	// Create backup with timestamp
	backupName := fmt.Sprintf("auto-backup-%s.json", time.Now().Format("2006-01-02-150405"))
	backupPath := filepath.Join(backupDir, backupName)
	
	if err := copyFile(dbPath, backupPath); err != nil {
		return fmt.Errorf("failed to create auto-backup: %w", err)
	}
	
	log.Printf("Auto-backup created: %s", backupPath)
	
	// Clean up old backups
	if err := cleanOldBackups(backupDir, 5); err != nil {
		log.Printf("Warning: failed to clean old backups: %v", err)
	}
	
	return nil
}

// Helper functions

func getDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, "Library", "Application Support", "btrack", "btrack-data.json")
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	return destFile.Sync()
}

func cleanOldBackups(backupDir string, keepCount int) error {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}
	
	type backupFile struct {
		name    string
		modTime time.Time
	}
	
	var backups []backupFile
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
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
