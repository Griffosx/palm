package config

import (
	"fmt"
	"math/rand"
	"palm/src/entities"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// PalmDB creates a new database connection.
// If inMemory is true, an in-memory SQLite database is created.
// The identifier parameter can be used to create separate in-memory databases
// with consistent names (useful for tests).
// When inMemory is false, identifier can be used to specify a custom file path.
func PalmDB(inMemory bool, identifier ...string) (*gorm.DB, error) {
	dbPath := "palm.sqlite"
	if inMemory {
		dbName := "memdb1"
		if len(identifier) > 0 && identifier[0] != "" {
			// Use the provided identifier for a unique but consistent database name
			dbName = fmt.Sprintf("memdb_%s", identifier[0])
		}
		dbPath = fmt.Sprintf("file:%s?mode=memory&cache=shared", dbName)
	} else if len(identifier) > 0 && identifier[0] != "" {
		// Use the provided identifier as a custom file path
		dbPath = identifier[0]
	}

	fmt.Printf("\n\nUsing database path: %s\n\n\n", dbPath)

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Enable foreign key constraints
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	// Run migrations for all models
	err = db.AutoMigrate(
		&entities.Account{},
		&entities.Message{},
		&entities.Recipient{},
		&entities.Attachment{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
