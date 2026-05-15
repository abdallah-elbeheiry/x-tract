package models

import "github.com/google/uuid"

// Group represents a logical group for guest employees.
type Group struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NewGroup is the request payload for creating a group.
type NewGroup struct {
	Name string `json:"name" binding:"required"`
}

// UpdateGroup is the request payload for updating a group.
type UpdateGroup struct {
	Name *string `json:"name"`
}
