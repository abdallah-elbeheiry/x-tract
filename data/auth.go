package data

import (
	"context"
	"database/sql"
	"errors"
	"x-tract/models"
)

// AuthStore handles credential validation against the users table.
type AuthStore struct {
	db *Database
}

// NewAuthStore builds an AuthStore backed by the shared database.
func NewAuthStore(db *Database) *AuthStore {
	return &AuthStore{db: db}
}

// Authenticate verifies a user's email and password and returns the public user record.
func (s *AuthStore) Authenticate(ctx context.Context, email string, password string) (*models.User, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	var user models.User
	var number sql.NullString

	row := s.db.conn.QueryRowContext(ctx, `
		SELECT id, name, email, number, role, hash, salt
		FROM users
		WHERE email = $1
	`, email)

	if err := row.Scan(&user.ID, &user.Username, &user.Email, &number, &user.Role, &user.Hash, &user.Salt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, normalizeDBError(err)
	}

	user.Number = nullStringPtr(number)
	if !VerifyPassword(password, user.Hash, user.Salt) {
		return nil, ErrInvalidCredentials
	}

	user.Hash = nil
	user.Salt = nil
	return &user, nil
}
