package data

import (
	"context"
	"database/sql"
	"log"
	"time"
	"x-tract/models"

	"github.com/google/uuid"
)

type Database struct {
	conn *sql.DB
}

// CreateUser handles the two-step "Is-A" insertion using a Transaction
func (d *Database) CreateUser(ctx context.Context, u *models.User) error {
	// 1. Create a context with a timeout for concurrency safety
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 2. Start a transaction (Atomic operation)
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			// Log the error, but don't return it since we want to return the original error if there was one
			// Use a custom logger
			log.Printf("transaction rollback error: %v", err)
			return
		}
	}(tx)

	query := `INSERT INTO users (id, name, email, hash, salt, number, role)
             VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.ExecContext(ctx, query, u.ID, u.Username, u.Email, u.Hash, u.Salt, u.Number, u.Role)
	if err != nil {
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// GetUserByID demonstrates a high-speed read operation with a context timeout for concurrency safety
func (d *Database) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	u := &models.User{}
	query := `SELECT id, name, email, number, role FROM users WHERE id = $1`

	err := d.conn.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.Number, &u.Role,
	)

	if err != nil {
		return nil, err
	}
	return u, nil
}
