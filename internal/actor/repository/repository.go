package repository

import "intern/models"

type ActorRepositoryI interface {
	Create(a *models.Actor) error
	Get(id int) (*models.Actor, error)
	Update(a *models.Actor) error
	Delete(id int) error
	GetMoviesByActor(id int) ([]models.Movie, error)
}
