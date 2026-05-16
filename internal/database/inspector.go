package database

import (
	"database/sql"
	"fmt"
	"tracking/tools"
)

func PrintDatabaseStructure(db *sql.DB) error {

	tablesQuery := `
	SELECT table_name
	FROM information_schema.tables
	WHERE table_schema = 'public'
	ORDER BY table_name;
	`

	rows, err := db.Query(tablesQuery)
	if err != nil {
		return fmt.Errorf(
			"failed to fetch tables: %w",
			err,
		)
	}

	defer rows.Close()

	fmt.Println(tools.Bold, tools.Cyan, "========== TABLES ==========", tools.Reset)

	for rows.Next() {

		var table string

		if err := rows.Scan(&table); err != nil {
			return err
		}

		fmt.Printf("\nTable: %s\n", table)

		if err := printColumns(db, table); err != nil {
			return err
		}

		if err := printIndexes(db, table); err != nil {
			return err
		}
	}

	return nil
}

func printColumns(
	db *sql.DB,
	table string,
) error {

	query := `
	SELECT
		column_name,
		data_type,
		is_nullable
	FROM information_schema.columns
	WHERE table_name = $1
	ORDER BY ordinal_position;
	`

	rows, err := db.Query(query, table)
	if err != nil {
		return err
	}

	defer rows.Close()

	fmt.Println("  Columns:")

	for rows.Next() {

		var name string
		var dataType string
		var nullable string

		if err := rows.Scan(
			&name,
			&dataType,
			&nullable,
		); err != nil {
			return err
		}

		fmt.Printf(
			"    - %s (%s) nullable=%s\n",
			tools.Bold, tools.Yellow, name, tools.Reset,
			tools.Bold, tools.Cyan, dataType, tools.Reset,
			tools.Bold, tools.Green, nullable, tools.Reset,
		)
	}

	return nil
}

func printIndexes(
	db *sql.DB,
	table string,
) error {

	query := `
	SELECT
		indexname,
		indexdef
	FROM pg_indexes
	WHERE tablename = $1;
	`

	rows, err := db.Query(query, table)
	if err != nil {
		return err
	}

	defer rows.Close()

	fmt.Println("  Indexes:")

	for rows.Next() {

		var name string
		var definition string

		if err := rows.Scan(
			&name,
			&definition,
		); err != nil {
			return err
		}

		fmt.Printf(
			"    - %s\n",
			tools.Bold, tools.Yellow, name, tools.Reset,
		)

		fmt.Printf(
			"      %s\n",
			tools.Bold, tools.White, definition, tools.Reset,
		)
	}

	return nil
}
