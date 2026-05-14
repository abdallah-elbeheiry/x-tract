package models

import "github.com/google/uuid"

type Role string

const (
	Role_Admin         Role = "ADMIN"
	Role_Customer      Role = "CUSTOMER"
	Role_Salesman      Role = "SALES_MAN"
	Role_GuestEmployee Role = "GUEST_EMPLOYEE"
)

func IsValidRole(r string) bool {
	role := Role(r)
	switch role {
	case Role_Admin, Role_Customer, Role_Salesman, Role_GuestEmployee:
		return true
	default:
		return false
	}
}

// ParseRole converts a string into a validated role.
func ParseRole(r string) (Role, bool) {
	role := Role(r)
	if !IsValidRole(r) {
		return "", false
	}
	return role, true
}

// User is the main identity record stored in the system.
type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Number   *string   `json:"number,omitempty"` // optional field
	Role     Role      `json:"role"`
	Hash     []byte    `json:"-"`
	Salt     []byte    `json:"-"`
}

// Is reports whether the user has the exact provided role.
func (u *User) Is(r Role) bool {
	return u.Role == r
}

// GetAuthInfo returns the stored password hash and salt.
func (u *User) GetAuthInfo() ([]byte, []byte) {
	return u.Hash, u.Salt
}

// SetAuthInfo stores the password hash and salt on the user.
func (u *User) SetAuthInfo(hash, salt []byte) {
	u.Hash = hash
	u.Salt = salt
}

// HasRole reports whether the user matches any of the provided roles.
func (u *User) HasRole(roles ...Role) bool {
	for _, r := range roles {
		if u.Role == r {
			return true
		}
	}
	return false
}

//logic added later
