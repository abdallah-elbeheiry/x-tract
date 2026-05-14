package models

// NewUserFields contains the shared fields required to create a typed user.
type NewUserFields struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Number   string `json:"number" binding:"omitempty,e164"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateUserFields contains the shared fields allowed when updating a typed user.
type UpdateUserFields struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" binding:"omitempty,email"`
	Number   *string `json:"number" binding:"omitempty,e164"`
	Password *string `json:"password" binding:"omitempty,min=8"`
}

// NewAdmin is the request payload for creating an admin.
type NewAdmin struct {
	NewUserFields
}

// UpdateAdmin is the request payload for updating an admin.
type UpdateAdmin struct {
	UpdateUserFields
}

// NewCustomer is the request payload for creating a customer.
type NewCustomer struct {
	NewUserFields
	Stats string `json:"stats"`
}

// UpdateCustomer is the partial request payload for updating a customer.
type UpdateCustomer struct {
	UpdateUserFields
	Stats *string `json:"stats"`
}

// NewSalesman is the request payload for creating a salesman.
type NewSalesman struct {
	NewUserFields
	Stats string `json:"stats"`
}

// UpdateSalesman is the partial request payload for updating a salesman.
type UpdateSalesman struct {
	UpdateUserFields
	Stats *string `json:"stats"`
}

// NewGuestEmployee is the request payload for creating a guest employee.
type NewGuestEmployee struct {
	NewUserFields
	GroupID *string `json:"group_id" binding:"omitempty,uuid"`
}

// UpdateGuestEmployee is the partial request payload for updating a guest employee.
type UpdateGuestEmployee struct {
	UpdateUserFields
	GroupID *string `json:"group_id" binding:"omitempty,uuid"`
}
