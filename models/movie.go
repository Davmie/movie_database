package models

import "time"

type Movie struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	ReleaseDate time.Time `json:"releaseDate" db:"releaseDate"`
	Rating      int       `json:"rating" db:"rating"`
}
