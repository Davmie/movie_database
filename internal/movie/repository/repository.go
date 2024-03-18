package repository

import "intern/models"

type MovieRepositoryI interface {
	Create(m *models.Movie) error
	Get(id int) (*models.Movie, error)
	Update(m *models.Movie) error
	Delete(id int) error
	GetMoviesSorted(sortingColumn string) ([]models.Movie, error)
	GetActorsByMovie(id int) ([]models.Actor, error)
	GetMoviesByTitle(title string) ([]models.Movie, error)
}
