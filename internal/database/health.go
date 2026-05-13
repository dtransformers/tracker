package database

import (
	"context"
	"database/sql"
	"time"
)

func Health(db *sql.DB) error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer cancel()

	return db.PingContext(ctx)
}
