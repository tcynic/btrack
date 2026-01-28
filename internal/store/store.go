package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Store holds the application state in memory and manages persistence
type Store struct {
	data     *Data
	filePath string
	mu       sync.RWMutex
}

// New creates a new Store instance with a custom file path
func New(filePath string) *Store {
	return &Store{
		data:     NewData(),
		filePath: filePath,
	}
}

// NewStore creates a new Store instance with the default file path
func NewStore() *Store {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir not available
		return New("./btrack-data.json")
	}
	filePath := filepath.Join(homeDir, "Library", "Application Support", "btrack", "btrack-data.json")
	return New(filePath)
}

// Load reads the JSON file from disk into memory
// If the file doesn't exist, initializes with empty data
func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Check if file exists
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		// File doesn't exist, initialize with empty data
		s.data = NewData()
		// Save initial empty state
		return s.saveUnlocked()
	}

	// Read file
	fileData, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var data Data
	if err := json.Unmarshal(fileData, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	s.data = &data
	return nil
}

// Save writes the current state to disk atomically
// Uses write-to-temp-then-rename pattern for safety
func (s *Store) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveUnlocked()
}

// saveUnlocked performs the actual save without locking (internal use)
func (s *Store) saveUnlocked() error {
	// Marshal to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to temporary file first
	tempPath := s.filePath + ".tmp"
	if err := os.WriteFile(tempPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Atomically rename temp file to actual file
	if err := os.Rename(tempPath, s.filePath); err != nil {
		// Clean up temp file on error
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// NextID generates and returns the next ID for an entity type
func (s *Store) NextID(entityType string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.data.NextIDs[entityType]
	s.data.NextIDs[entityType] = id + 1
	return id
}

// GetFilePath returns the path to the data file
func (s *Store) GetFilePath() string {
	return s.filePath
}
