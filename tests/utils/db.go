package utils

import (
	"os"
	"palm/src/config"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestDBPath returns a unique file path for a test database
func TestDBPath(testName string) string {
	// Get the project root directory
	rootDir, err := os.Getwd()
	for !isRootDir(rootDir) && err == nil {
		rootDir = filepath.Dir(rootDir)
	}

	dbDir := filepath.Join(rootDir, "tests", "temp_db")
	// Ensure the directory exists
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		// If we can't create the directory, fall back to temp dir
		dbDir = os.TempDir()
	}
	return filepath.Join(dbDir, "palm_test_"+testName+".sqlite")
}

// isRootDir checks if the given directory is the project root
// by looking for common project files
func isRootDir(dir string) bool {
	rootIndicators := []string{"go.mod", ".git"}
	for _, indicator := range rootIndicators {
		if _, err := os.Stat(filepath.Join(dir, indicator)); err == nil {
			return true
		}
	}
	return false
}

// SetupTestDB creates a file-based SQLite database for testing
// It also runs migrations for all entities
func SetupTestDB(t *testing.T) *gorm.DB {
	// Create a unique file path for this test
	dbPath := TestDBPath(t.Name())

	// Use file-based database (inMemory = false)
	db, err := config.PalmDB(false, dbPath)
	require.NoError(t, err)
	return db
}

// TeardownTestDB removes the test database file after test completion
func TeardownTestDB(t *testing.T) {
	dbPath := TestDBPath(t.Name())
	err := os.Remove(dbPath)
	if err != nil && !os.IsNotExist(err) {
		t.Logf("Failed to remove test database file %s: %v", dbPath, err)
	}
}
