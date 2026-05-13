package models

import "github.com/google/uuid"

type Customer struct {
	UserID          uuid.UUID `json:"user_id"` // PK, same as User
	OtherStatistics string    `json:"other_statistics"`
}

type SalesMan struct {
	UserID          uuid.UUID `json:"user_id"` // PK, same as User
	OtherStatistics string    `json:"other_statistics"`
}

type GuestEmployee struct {
	UserID  uuid.UUID `json:"user_id"`  // PK, same as User
	GroupID uuid.UUID `json:"group_id"` // FK to Group
}

//logic added later
