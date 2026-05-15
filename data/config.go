package data

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	URL          string
	MaxOpenConns int
	MaxIdleConns int
}

// NewPostgres creates a tuned connection pool
func NewPostgres(cfg PostgresConfig) (*Database, error) {
	if cfg.URL == "" {
		return nil, errors.New("database url is required")
	}

	db, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Performance Tuning
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return &Database{conn: db}, nil
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

// RunMigrations applies embedded SQL migrations to the configured database.
func RunMigrations(db *Database) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db.conn, "migrations"); err != nil {
		return err
	}
	return nil
}
