package models

import (
	"time"

	"github.com/google/uuid"
)

// Operation represents a completed collection or sales operation.
type Operation struct {
	ID         uuid.UUID `json:"id"`
	SalesmanID uuid.UUID `json:"salesman_id"` // FK
	CustomerID uuid.UUID `json:"customer_id"` // FK, customer can be group or customer
	Location   string    `json:"location"`    // Google Maps API string
	FilmWeight float64   `json:"film_weight"` // number/weight of films
	TimeDate   time.Time `json:"time_date"`
}

//logic added later
