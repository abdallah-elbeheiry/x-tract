package data

import (
	"context"
	"database/sql"
	"x-tract/models"

	"github.com/google/uuid"
)

const customerSelectColumns = `
	u.id, u.name, u.email, u.number, u.role, c.stats
`

type CustomerStore struct {
	db *Database
}

// NewCustomerStore builds a CustomerStore backed by the shared database.
func NewCustomerStore(db *Database) *CustomerStore {
	return &CustomerStore{db: db}
}

// List returns all customers.
func (s *CustomerStore) List(ctx context.Context) ([]models.Customer, error) {
	return s.db.listCustomers(ctx)
}

// Create inserts a new customer.
func (s *CustomerStore) Create(ctx context.Context, input *models.NewCustomer) (*models.Customer, error) {
	return s.db.createCustomer(ctx, input)
}

// GetByID loads a single customer.
func (s *CustomerStore) GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	return s.db.getCustomerByID(ctx, id)
}

// Update applies a partial update to a customer.
func (s *CustomerStore) Update(ctx context.Context, id uuid.UUID, input *models.UpdateCustomer) (*models.Customer, error) {
	return s.db.updateCustomer(ctx, id, input)
}

// Delete removes a customer.
func (s *CustomerStore) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.deleteCustomer(ctx, id)
}

func (d *Database) createCustomer(ctx context.Context, input *models.NewCustomer) (*models.Customer, error) {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	tx, err := d.conn.BeginTx(ctx, nil)
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
		Role:     models.Role_Customer,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO customer (user_id, stats)
		VALUES ($1, $2)
	`, userID, nullableString(input.Stats)); err != nil {
		return nil, normalizeDBError(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, normalizeDBError(err)
	}

	return d.getCustomerByID(context.Background(), userID)
}

// listCustomers returns customers joined with their user records.
func (d *Database) listCustomers(ctx context.Context) ([]models.Customer, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	rows, err := d.conn.QueryContext(ctx, `
		SELECT `+customerSelectColumns+`
		FROM customer c
		JOIN users u ON u.id = c.user_id
		ORDER BY u.name, u.id
	`)
	if err != nil {
		return nil, normalizeDBError(err)
	}
	defer rows.Close()

	return collectRows(rows, scanCustomer)
}

// getCustomerByID loads one customer together with its user data.
func (d *Database) getCustomerByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	row := d.conn.QueryRowContext(ctx, `
		SELECT `+customerSelectColumns+`
		FROM customer c
		JOIN users u ON u.id = c.user_id
		WHERE c.user_id = $1
	`, id)

	return scanCustomer(row)
}

// updateCustomer updates the customer row and the linked user row inside one
// transaction, then returns the fresh joined view.
func (d *Database) updateCustomer(ctx context.Context, id uuid.UUID, input *models.UpdateCustomer) (*models.Customer, error) {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := ensureCustomerExists(ctx, tx, id); err != nil {
		return nil, err
	}

	if err := updateUserWithExecutor(ctx, tx, id, input.UpdateUserFields); err != nil {
		return nil, err
	}

	if input.Stats != nil {
		if err := execUpdate(ctx, tx, `UPDATE customer SET stats = $1 WHERE user_id = $2`, nullableString(*input.Stats), id); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, normalizeDBError(err)
	}

	return d.getCustomerByID(context.Background(), id)
}

// deleteCustomer removes the customer by deleting the parent user row.
func (d *Database) deleteCustomer(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	return execDelete(ctx, d.conn, `DELETE FROM users WHERE id = $1 AND role = $2`, id, models.Role_Customer)
}

func ensureCustomerExists(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM customer WHERE user_id = $1)`, id).Scan(&exists); err != nil {
		return normalizeDBError(err)
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}
