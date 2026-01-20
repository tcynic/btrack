package main

import (
	"btrack/internal/database"
	"btrack/internal/services/system"
)

// BackupInfo provides information about the current database
type BackupInfo = system.BackupInfo

// GetBackupInfo returns information about the current database
func (a *App) GetBackupInfo() (*BackupInfo, error) {
	return a.systemService.GetBackupInfo()
}

// CreateBackup creates a backup of the database to a user-selected location
func (a *App) CreateBackup() (string, error) {
	return a.systemService.CreateBackup()
}

// RestoreBackup restores the database from a backup file
func (a *App) RestoreBackup() error {
	return a.systemService.RestoreBackup(
		func() error { return a.db.Close() },
		func() error {
			db, err := database.Initialize()
			if err != nil {
				return err
			}
			a.db = db
			return nil
		},
	)
}

// autoBackup creates an automatic backup before destructive operations
func (a *App) autoBackup() error {
	return a.systemService.AutoBackup()
}
