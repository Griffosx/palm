package config

import (
	"palm/src/entities"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func PalmDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("palm.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = db.AutoMigrate(&entities.Account{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
