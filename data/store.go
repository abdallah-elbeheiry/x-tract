package data

import (
	"context"
	"database/sql"
	"time"
)

const (
	defaultWriteTimeout = 5 * time.Second
	defaultReadTimeout  = 2 * time.Second
)

// Database provides persistence operations for the server.
//
// The type intentionally stays small: connection ownership lives here while
// resource-specific CRUD is split across focused files in the same package.
type Database struct {
	conn *sql.DB
}

// Close releases the underlying database connection pool.
func (d *Database) Close() error {
	if d == nil || d.conn == nil {
		return nil
	}
	return d.conn.Close()
}

func withReadTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultReadTimeout)
}

func withWriteTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, defaultWriteTimeout)
}
