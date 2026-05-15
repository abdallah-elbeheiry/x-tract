package data

import (
	"context"
	"fmt"
	"strings"
	"x-tract/models"

	"github.com/google/uuid"
)

// GroupStore exposes CRUD behavior for groups.
type GroupStore struct {
	db *Database
}

// NewGroupStore builds a GroupStore backed by the shared database.
func NewGroupStore(db *Database) *GroupStore {
	return &GroupStore{db: db}
}

// List returns all groups.
func (s *GroupStore) List(ctx context.Context) ([]models.Group, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	rows, err := s.db.conn.QueryContext(ctx, `SELECT id, name FROM groups ORDER BY name, id`)
	if err != nil {
		return nil, normalizeDBError(err)
	}
	defer rows.Close()

	return collectRows(rows, scanGroup)
}

// Create inserts a new group.
func (s *GroupStore) Create(ctx context.Context, input *models.NewGroup) (*models.Group, error) {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	group := &models.Group{
		ID:   uuid.New(),
		Name: input.Name,
	}

	if _, err := s.db.conn.ExecContext(ctx, `INSERT INTO groups (id, name) VALUES ($1, $2)`, group.ID, group.Name); err != nil {
		return nil, normalizeDBError(err)
	}

	return group, nil
}

// GetByID loads a single group.
func (s *GroupStore) GetByID(ctx context.Context, id uuid.UUID) (*models.Group, error) {
	ctx, cancel := withReadTimeout(ctx)
	defer cancel()

	row := s.db.conn.QueryRowContext(ctx, `SELECT id, name FROM groups WHERE id = $1`, id)
	return scanGroup(row)
}

// Update applies a partial update to a group.
func (s *GroupStore) Update(ctx context.Context, id uuid.UUID, input *models.UpdateGroup) (*models.Group, error) {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	setClauses := make([]string, 0, 1)
	args := make([]any, 0, 2)
	nextArg := 1

	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", nextArg))
		args = append(args, *input.Name)
		nextArg++
	}

	if len(setClauses) == 0 {
		return s.GetByID(ctx, id)
	}

	query := fmt.Sprintf(`UPDATE groups SET %s WHERE id = $%d RETURNING id, name`, strings.Join(setClauses, ", "), nextArg)
	row := s.db.conn.QueryRowContext(ctx, query, append(args, id)...)
	return scanGroup(row)
}

// Delete removes a group.
func (s *GroupStore) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := withWriteTimeout(ctx)
	defer cancel()

	return execDelete(ctx, s.db.conn, `DELETE FROM groups WHERE id = $1`, id)
}
