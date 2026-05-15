package data

import (
	"context"
	"database/sql"
	"x-tract/models"

	"github.com/google/uuid"
)

const guestEmployeeSelectColumns = `
	u.id, u.name, u.email, u.number, u.role, ge.group_id::text, grp.name
`

// GuestEmployeeStore exposes CRUD behavior for guest employee resources.
type GuestEmployeeStore struct {
	db *Database
}

// NewGuestEmployeeStore builds a GuestEmployeeStore backed by the shared database.
func NewGuestEmployeeStore(db *Database) *GuestEmployeeStore {
	return &GuestEmployeeStore{db: db}
}

// List returns all guest employees.
func (s *GuestEmployeeStore) List(ctx context.Context) ([]models.GuestEmployee, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	rows, err := s.db.conn.QueryContext(ctx, `
		SELECT `+guestEmployeeSelectColumns+`
		FROM guest_employee ge
		JOIN users u ON u.id = ge.user_id
		LEFT JOIN groups grp ON grp.id = ge.group_id
		ORDER BY u.name, u.id
	`)
	if err != nil {
		return nil, normalizeDBError(err)
	}
	defer rows.Close()

	return collectRows(rows, scanGuestEmployee)
}

// Create inserts a new guest employee.
func (s *GuestEmployeeStore) Create(ctx context.Context, input *models.NewGuestEmployee) (*models.GuestEmployee, error) {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	tx, err := s.db.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	userID, err := insertUser(ctx, tx, newUserInsert{
		Name:     input.Name,
		Email:    input.Email,
		Number:   input.Number,
		Role:     models.Role_GuestEmployee,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO guest_employee (user_id, group_id) VALUES ($1, $2)`, userID, nullableUUIDString(input.GroupID)); err != nil {
		return nil, normalizeDBError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, normalizeDBError(err)
	}

	return s.GetByID(context.Background(), userID)
}

// GetByID loads a single guest employee.
func (s *GuestEmployeeStore) GetByID(ctx context.Context, id uuid.UUID) (*models.GuestEmployee, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	row := s.db.conn.QueryRowContext(ctx, `
		SELECT `+guestEmployeeSelectColumns+`
		FROM guest_employee ge
		JOIN users u ON u.id = ge.user_id
		LEFT JOIN groups grp ON grp.id = ge.group_id
		WHERE ge.user_id = $1
	`, id)

	return scanGuestEmployee(row)
}

// Update applies a partial update to a guest employee.
func (s *GuestEmployeeStore) Update(ctx context.Context, id uuid.UUID, input *models.UpdateGuestEmployee) (*models.GuestEmployee, error) {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	tx, err := s.db.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := ensureGuestEmployeeExists(ctx, tx, id); err != nil {
		return nil, err
	}

	if err := updateUserWithExecutor(ctx, tx, id, input.UpdateUserFields); err != nil {
		return nil, err
	}

	if input.GroupID != nil {
		if err := execUpdate(ctx, tx, `UPDATE guest_employee SET group_id = $1 WHERE user_id = $2`, nullableUUIDString(input.GroupID), id); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, normalizeDBError(err)
	}

	return s.GetByID(context.Background(), id)
}

// Delete removes a guest employee.
func (s *GuestEmployeeStore) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	return execDelete(ctx, s.db.conn, `DELETE FROM users WHERE id = $1 AND role = $2`, id, models.Role_GuestEmployee)
}

func ensureGuestEmployeeExists(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM guest_employee WHERE user_id = $1)`, id).Scan(&exists); err != nil {
		return normalizeDBError(err)
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}

func nullableUUIDString(value *string) any {
	if value == nil || *value == "" {
		return nil
	}
	return *value
}
