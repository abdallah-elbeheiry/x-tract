package data

import (
	"context"
	"database/sql"
	"x-tract/models"

	"github.com/google/uuid"
)

const salesmanSelectColumns = `
	u.id, u.name, u.email, u.number, u.role, s.stats
`

// SalesmanStore exposes CRUD behavior for salesman resources.
type SalesmanStore struct {
	db *Database
}

// NewSalesmanStore builds a SalesmanStore backed by the shared database.
func NewSalesmanStore(db *Database) *SalesmanStore {
	return &SalesmanStore{db: db}
}

// List returns all salesmen.
func (s *SalesmanStore) List(ctx context.Context) ([]models.Salesman, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	rows, err := s.db.conn.QueryContext(ctx, `
		SELECT `+salesmanSelectColumns+`
		FROM salesman s
		JOIN users u ON u.id = s.user_id
		ORDER BY u.name, u.id
	`)
	if err != nil {
		return nil, normalizeDBError(err)
	}
	defer rows.Close()

	return collectRows(rows, scanSalesman)
}

// Create inserts a new salesman.
func (s *SalesmanStore) Create(ctx context.Context, input *models.NewSalesman) (*models.Salesman, error) {
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
		Role:     models.Role_Salesman,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO salesman (user_id, stats) VALUES ($1, $2)`, userID, nullableString(input.Stats)); err != nil {
		return nil, normalizeDBError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, normalizeDBError(err)
	}

	return s.GetByID(context.Background(), userID)
}

// GetByID loads a single salesman.
func (s *SalesmanStore) GetByID(ctx context.Context, id uuid.UUID) (*models.Salesman, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	row := s.db.conn.QueryRowContext(ctx, `
		SELECT `+salesmanSelectColumns+`
		FROM salesman s
		JOIN users u ON u.id = s.user_id
		WHERE s.user_id = $1
	`, id)

	return scanSalesman(row)
}

// Update applies a partial update to a salesman.
func (s *SalesmanStore) Update(ctx context.Context, id uuid.UUID, input *models.UpdateSalesman) (*models.Salesman, error) {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	tx, err := s.db.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := ensureSalesmanExists(ctx, tx, id); err != nil {
		return nil, err
	}

	if err := updateUserWithExecutor(ctx, tx, id, input.UpdateUserFields); err != nil {
		return nil, err
	}

	if input.Stats != nil {
		if err := execUpdate(ctx, tx, `UPDATE salesman SET stats = $1 WHERE user_id = $2`, nullableString(*input.Stats), id); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, normalizeDBError(err)
	}

	return s.GetByID(context.Background(), id)
}

// Delete removes a salesman.
func (s *SalesmanStore) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	return execDelete(ctx, s.db.conn, `DELETE FROM users WHERE id = $1 AND role = $2`, id, models.Role_Salesman)
}

func ensureSalesmanExists(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM salesman WHERE user_id = $1)`, id).Scan(&exists); err != nil {
		return normalizeDBError(err)
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}
