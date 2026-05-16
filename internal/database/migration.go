package database

import (
	"context"
	"fmt"
	"log"
	"tracking/internal/ent"

	"entgo.io/ent/dialect/sql/schema"
)

func RunMigrations(
	ctx context.Context,
	client *ent.Client,
) error {

	log.Println("Starting database migrations...")

	err := client.Schema.Create(
		ctx,

		schema.WithForeignKeys(true),

		schema.WithGlobalUniqueID(false),
	)

	if err != nil {
		return fmt.Errorf(
			"failed to create schema: %w",
			err,
		)
	}

	log.Println("Database migrations completed")

	return nil
}
