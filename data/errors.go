package data

import "errors"

var (
	// ErrNotFound reports that a requested row does not exist.
	ErrNotFound = errors.New("resource not found")
	// ErrConflict reports a uniqueness or duplicate-key violation.
	ErrConflict = errors.New("resource conflict")
	// ErrInvalidRole reports an unsupported user role value.
	ErrInvalidRole = errors.New("invalid role")
	// ErrForeignKeyViolation reports an invalid referenced row.
	ErrForeignKeyViolation = errors.New("invalid reference")
	// ErrInvalidCredentials reports a login failure.
	ErrInvalidCredentials = errors.New("invalid credentials")
)
