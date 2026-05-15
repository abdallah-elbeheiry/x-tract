package data

import (
	"context"
	"x-tract/models"

	"github.com/google/uuid"
)

const roleUserSelectColumns = `id, name, email, number, role`

// AdminStore exposes CRUD behavior for admin resources.
type AdminStore struct {
	db *Database
}

// NewAdminStore builds an AdminStore backed by the shared database.
func NewAdminStore(db *Database) *AdminStore {
	return &AdminStore{db: db}
}

// List returns all admins.
func (s *AdminStore) List(ctx context.Context) ([]models.Admin, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	rows, err := s.db.conn.QueryContext(ctx, `SELECT `+roleUserSelectColumns+` FROM users WHERE role = $1 ORDER BY name, id`, models.Role_Admin)
	if err != nil {
		return nil, normalizeDBError(err)
	}
	defer rows.Close()

	return collectRows(rows, scanRoleUser)
}

// Create inserts a new admin.
func (s *AdminStore) Create(ctx context.Context, input *models.NewAdmin) (*models.Admin, error) {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	userID, err := insertUser(ctx, s.db.conn, newUserInsert{
		Name:     input.Name,
		Email:    input.Email,
		Number:   input.Number,
		Role:     models.Role_Admin,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}

	return s.GetByID(context.Background(), userID)
}

// GetByID loads a single admin.
func (s *AdminStore) GetByID(ctx context.Context, id uuid.UUID) (*models.Admin, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	row := s.db.conn.QueryRowContext(ctx, `SELECT `+roleUserSelectColumns+` FROM users WHERE id = $1 AND role = $2`, id, models.Role_Admin)
	return scanRoleUser(row)
}

// Update applies a partial update to an admin.
func (s *AdminStore) Update(ctx context.Context, id uuid.UUID, input *models.UpdateAdmin) (*models.Admin, error) {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	if err := updateRoleWithExecutor(ctx, s.db.conn, id, models.Role_Admin, input.UpdateUserFields); err != nil {
		return nil, err
	}

	if emptyUserUpdate(input.UpdateUserFields) {
		return s.GetByID(ctx, id)
	}
	return s.GetByID(context.Background(), id)
}

// Delete removes an admin.
func (s *AdminStore) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	return execDelete(ctx, s.db.conn, `DELETE FROM users WHERE id = $1 AND role = $2`, id, models.Role_Admin)
}

func emptyUserUpdate(input models.UpdateUserFields) bool {
	return input.Name == nil && input.Email == nil && input.Number == nil && input.Password == nil
}
