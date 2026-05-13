package database

import "time"

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration

	EnableDebug bool
}
