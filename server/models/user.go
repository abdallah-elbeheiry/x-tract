package models

type Role string

const (
	Admin         Role = "ADMIN"
	Customer      Role = "CUSTOMER"
	SalesMan      Role = "SALES_MAN"
	GuestEmployee Role = "GUEST_EMPLOYEE"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
	Hash     []byte `json:"hash"`
	Salt     []byte `json:"salt"`
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

func NewUser(id int, username, email string, role Role) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
	}
}
