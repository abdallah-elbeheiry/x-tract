package models

import "github.com/google/uuid"

type Role string

const (
	Role_Admin         Role = "ADMIN"
	Role_Customer      Role = "CUSTOMER"
	Role_SalesMan      Role = "SALES_MAN"
	Role_GuestEmployee Role = "GUEST_EMPLOYEE"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Number   *string   `json:"number,omitempty"` // optional field
	Role     Role      `json:"role"`
	Hash     []byte    `json:"-"`
	Salt     []byte    `json:"-"`
}

func (u *User) Is(r Role) bool {
	return u.Role == r
}

func (u *User) GetAuthInfo() ([]byte, []byte) {
	return u.Hash, u.Salt
}

func (u *User) SetAuthInfo(hash, salt []byte) {
	u.Hash = hash
	u.Salt = salt
}

func (u *User) HasRole(roles ...Role) bool {
	for _, r := range roles {
		if u.Role == r {
			return true
		}
	}
	return false
}

func NewUser(id uuid.UUID, username, email string, role Role) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
	}
}

//logic added later
