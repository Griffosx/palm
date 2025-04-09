package utils_test

import (
	"os"
	"palm/tests/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseFile(t *testing.T) {
	// Get the database path that will be used
	dbPath := utils.TestDBPath(t.Name())

	// Make sure the file doesn't exist before the test
	_, err := os.Stat(dbPath)
	assert.True(t, os.IsNotExist(err), "Database file should not exist before test")

	// Setup the database
	db := utils.SetupTestDB(t)
	require.NotNil(t, db, "Database should be created successfully")

	// Verify the file now exists
	_, err = os.Stat(dbPath)
	assert.NoError(t, err, "Database file should exist during test")

	// Clean up after the test
	utils.TeardownTestDB(t)

	// Verify the file was deleted
	_, err = os.Stat(dbPath)
	assert.True(t, os.IsNotExist(err), "Database file should be deleted after teardown")
}
