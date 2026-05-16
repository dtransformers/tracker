package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	"tracking/internal/ent"

	_ "github.com/lib/pq"
)

type Database struct {
	Client *ent.Client
	SQL    *sql.DB
}

func Connect(cfg Config) (*Database, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	driver := entsql.OpenDB(dialect.Postgres, sqlDB)

	client := ent.NewClient(ent.Driver(driver))

	if err := warmUp(sqlDB, cfg.MaxOpenConns); err != nil {
		return nil, err
	}

	fmt.Printf("Database Connected successfully %s:%d\n", cfg.Host, cfg.Port)
	defer main(dsn)
	return &Database{
		Client: client,
		SQL:    sqlDB,
	}, nil
}

func warmUp(db *sql.DB, maxOpenConnections int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	for range maxOpenConnections {
		conn, err := db.Conn(ctx)
		if err != nil {
			return fmt.Errorf("warm connection: %w", err)
		}
		_ = conn.PingContext(ctx)
		_ = conn.Close()
	}

	return nil
}
