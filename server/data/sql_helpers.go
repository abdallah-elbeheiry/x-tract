package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"x-tract/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type scanner interface {
	Scan(dest ...any) error
}

type execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type userPatch struct {
	setClauses []string
	args       []any
	nextArg    int
}

type newUserInsert struct {
	Name     string
	Email    string
	Number   string
	Role     models.Role
	Password string
}

func (p userPatch) empty() bool {
	return len(p.setClauses) == 0
}

func newUserPatch(input models.UpdateUserFields) (userPatch, error) {
	patch := userPatch{
		setClauses: make([]string, 0, 5),
		args:       make([]any, 0, 6),
		nextArg:    1,
	}

	addField := func(column string, value any) {
		patch.setClauses = append(patch.setClauses, fmt.Sprintf("%s = $%d", column, patch.nextArg))
		patch.args = append(patch.args, value)
		patch.nextArg++
	}

	if input.Name != nil {
		addField("name", *input.Name)
	}
	if input.Email != nil {
		addField("email", *input.Email)
	}
	if input.Number != nil {
		addField("number", nullableString(*input.Number))
	}
	if input.Password != nil {
		hash, salt, err := hashPassword(*input.Password)
		if err != nil {
			return userPatch{}, err
		}
		addField("hash", hash)
		addField("salt", salt)
	}

	return patch, nil
}

func updateUserWithExecutor(ctx context.Context, db execer, id uuid.UUID, input models.UpdateUserFields) error {
	patch, err := newUserPatch(input)
	if err != nil {
		return err
	}
	if patch.empty() {
		return nil
	}

	query := buildUpdateByIDQuery("users", patch.setClauses, patch.nextArg)
	args := append(patch.args, id)
	return execUpdate(ctx, db, query, args...)
}

func updateRoleWithExecutor(ctx context.Context, db execer, id uuid.UUID, role models.Role, input models.UpdateUserFields) error {
	patch, err := newUserPatch(input)
	if err != nil {
		return err
	}

	if patch.empty() {
		return nil
	}

	query := buildUpdateByIDAndRoleQuery("users", patch.setClauses, patch.nextArg, patch.nextArg+1)
	args := append(patch.args, id, role)
	return execUpdate(ctx, db, query, args...)
}

func insertUser(ctx context.Context, db execer, input newUserInsert) (uuid.UUID, error) {
	userID := uuid.New()

	hash, salt, err := hashPassword(input.Password)
	if err != nil {
		return uuid.Nil, err
	}

	const query = `
		INSERT INTO users (id, name, email, hash, salt, number, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := db.ExecContext(ctx, query, userID, input.Name, input.Email, hash, salt, nullableString(input.Number), input.Role); err != nil {
		return uuid.Nil, normalizeDBError(err)
	}

	return userID, nil
}

func scanRoleUser(row scanner) (*models.Admin, error) {
	user, err := scanUser(row)
	if err != nil {
		return nil, err
	}
	return &models.Admin{User: user}, nil
}

func scanSalesman(row scanner) (*models.Salesman, error) {
	var user models.User
	var number sql.NullString
	var stats sql.NullString

	if err := row.Scan(&user.ID, &user.Username, &user.Email, &number, &user.Role, &stats); err != nil {
		return nil, normalizeDBError(err)
	}

	user.Number = nullStringPtr(number)
	return &models.Salesman{
		User:            &user,
		OtherStatistics: stats.String,
	}, nil
}

func scanGroup(row scanner) (*models.Group, error) {
	var group models.Group
	if err := row.Scan(&group.ID, &group.Name); err != nil {
		return nil, normalizeDBError(err)
	}
	return &group, nil
}

func execUpdate(ctx context.Context, db execer, query string, args ...any) error {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return normalizeDBError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func execDelete(ctx context.Context, db execer, query string, args ...any) error {
	return execUpdate(ctx, db, query, args...)
}

func collectRows[T any](rows *sql.Rows, scan func(scanner) (*T, error)) ([]T, error) {
	items := make([]T, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, normalizeDBError(err)
	}

	return items, nil
}

func scanUser(row scanner) (*models.User, error) {
	var user models.User
	var number sql.NullString

	if err := row.Scan(&user.ID, &user.Username, &user.Email, &number, &user.Role); err != nil {
		return nil, normalizeDBError(err)
	}

	user.Number = nullStringPtr(number)
	return &user, nil
}

func scanCustomer(row scanner) (*models.Customer, error) {
	var user models.User
	var number sql.NullString
	var stats sql.NullString

	if err := row.Scan(&user.ID, &user.Username, &user.Email, &number, &user.Role, &stats); err != nil {
		return nil, normalizeDBError(err)
	}

	user.Number = nullStringPtr(number)

	customer := &models.Customer{
		User:            &user,
		OtherStatistics: stats.String,
	}

	return customer, nil
}

func scanGuestEmployee(row scanner) (*models.GuestEmployee, error) {
	var user models.User
	var number sql.NullString
	var groupID sql.NullString
	var groupName sql.NullString

	if err := row.Scan(&user.ID, &user.Username, &user.Email, &number, &user.Role, &groupID, &groupName); err != nil {
		return nil, normalizeDBError(err)
	}

	user.Number = nullStringPtr(number)

	guest := &models.GuestEmployee{User: &user}
	if groupID.Valid {
		parsed, err := uuid.Parse(groupID.String)
		if err != nil {
			return nil, err
		}
		guest.Group = &models.Group{
			ID:   parsed,
			Name: groupName.String,
		}
	}

	return guest, nil
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func buildUpdateByIDQuery(table string, setClauses []string, idArg int) string {
	query := "UPDATE " + table + " SET "
	for i, clause := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += clause
	}
	query += fmt.Sprintf(" WHERE id = $%d", idArg)
	return query
}

func buildUpdateByIDAndRoleQuery(table string, setClauses []string, idArg int, roleArg int) string {
	return buildUpdateByIDQuery(table, setClauses, idArg) + fmt.Sprintf(" AND role = $%d", roleArg)
}

func normalizeDBError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return ErrConflict
		case "23503":
			return ErrForeignKeyViolation
		}
	}

	return err
}
