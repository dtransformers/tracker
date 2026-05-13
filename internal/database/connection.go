package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}

func Connect(cfg Config) (*Database, error) {

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s timezone=UTC",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)

	gormDB, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	fmt.Printf(" Database Connected successfully %s:%d \n", cfg.Host, cfg.Port)

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("extract sql db: %w", err)
	}

	// Pool configuration
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Pre-warm connections
	if err := warmUp(sqlDB); err != nil {
		return nil, err
	}

	return &Database{
		Gorm: gormDB,
		SQL:  sqlDB,
	}, nil
}

func warmUp(db *sql.DB) error {

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Ping verifies connectivity
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	// Open initial connections
	for i := 0; i < 5; i++ {
		conn, err := db.Conn(ctx)
		if err != nil {
			return fmt.Errorf("warm connection: %w", err)
		}

		if err := conn.PingContext(ctx); err != nil {
			conn.Close()
			return err
		}

		conn.Close()
	}

	return nil
}
