package models

import "time"

type Actor struct {
	ID        int       `json:"id" db:"id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Gender    byte      `json:"gender" db:"gender"`
	Birthday  time.Time `json:"birthday" db:"birthday"`
}
