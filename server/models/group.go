package models

import "github.com/google/uuid"

// Group represents a logical group for guest employees.
type Group struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

//logic added later
