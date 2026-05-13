package database

import (
	"fmt"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {

	// get all files in ./models folder

	err := db.AutoMigrate()

	if err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}
