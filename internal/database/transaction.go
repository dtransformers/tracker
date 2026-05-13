package database

import (
	"context"

	"gorm.io/gorm"
)

func WithTransaction(
	ctx context.Context,
	db *gorm.DB,
	fn func(tx *gorm.DB) error,
) error {

	tx := db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
