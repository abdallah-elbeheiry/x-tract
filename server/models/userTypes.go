package models

import "github.com/google/uuid"

// Admin represents an admin user.
type Admin struct {
	User *User `json:"user" binding:"required"`
}

// Customer links a user record to customer-specific fields.
type Customer struct {
	User            *User  `json:"user" binding:"required"`
	OtherStatistics string `json:"other_statistics"`
}

// Salesman links a user record to salesman-specific fields.
type Salesman struct {
	User            *User  `json:"user" binding:"required"`
	OtherStatistics string `json:"other_statistics"`
}

// GuestEmployee links a user record to guest employee fields.
type GuestEmployee struct {
	User    *User      `json:"user" binding:"required"`
	GroupID *uuid.UUID `json:"group_id,omitempty"`
}
