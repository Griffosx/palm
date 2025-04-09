package utils

import (
	"palm/src/config"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
// It also runs migrations for all entities
func SetupTestDB(t *testing.T) *gorm.DB {
	// Use in-memory database with test name as identifier for consistent naming
	db, err := config.PalmDB(true, t.Name())
	require.NoError(t, err)
	return db
}
