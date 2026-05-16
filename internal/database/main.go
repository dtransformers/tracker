package database

import (
	"context"
	"database/sql"
	"fmt"
	"tracking/internal/ent"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

func newEntClient(db *sql.DB) (*ent.Client, error) {

	driver := entsql.OpenDB("postgres", db)

	client := ent.NewClient(
		ent.Driver(driver),
	)

	return client, nil
}

func newSQLClient(connectionString string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open postgres connection: %w",
			err,
		)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(
			"failed to ping postgres: %w",
			err,
		)
	}
	return db, nil
}

func main(connectionString string) {
	fmt.Print("Setup Required DB Actions: \n")

	ctx := context.Background()

	db, err := newSQLClient(connectionString)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	client, err := newEntClient(db)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	if err := RunMigrations(ctx, client); err != nil {
		panic(err)
	}

	if err := PrintDatabaseStructure(db); err != nil {
		panic(err)
	}
}
