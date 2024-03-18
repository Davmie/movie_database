package models

type MovieActor struct {
	ID      int `json:"id" db:"id"`
	MovieID int `json:"movie_id" db:"movie_id"`
	ActorID int `json:"actor_id" db:"actor_id"`
}
