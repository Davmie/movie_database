package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"intern/internal/movie/repository"
	"intern/models"
	"intern/pkg/logger"
)

type pgMovieRepo struct {
	Logger logger.Logger
	DB     *gorm.DB
}

func New(logger logger.Logger, db *gorm.DB) repository.MovieRepositoryI {
	return &pgMovieRepo{
		Logger: logger,
		DB:     db,
	}
}

func (mr *pgMovieRepo) Create(m *models.Movie) error {
	tx := mr.DB.Create(m)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgMovieRepo.Create error")
	}

	return nil
}

func (mr *pgMovieRepo) Get(id int) (*models.Movie, error) {
	var m models.Movie
	tx := mr.DB.Where("id = ?", id).Take(&m)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgMovieRepo.Get error")
	}

	return &m, nil
}

func (mr *pgMovieRepo) Update(m *models.Movie) error {
	tx := mr.DB.Omit("id").Updates(m)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgMovieRepo.Update error")
	}

	return nil
}

func (mr *pgMovieRepo) Delete(id int) error {
	tx := mr.DB.Delete(&models.Movie{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgMovieRepo.Delete error")
	}

	return nil
}

func (mr *pgMovieRepo) GetMoviesSorted(sortingColumn string) ([]models.Movie, error) {
	var movies []models.Movie

	tx := mr.DB.Order(sortingColumn).Find(&movies)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgMovieRepo.GetMoviesSorted error")
	}

	return movies, nil
}

func (mr *pgMovieRepo) GetActorsByMovie(id int) ([]models.Actor, error) {
	var actorIDs []int

	tx := mr.DB.Table("movies_actors").Select("actor_id").Where("movie_id = ?", id).Find(&actorIDs)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgMovieRepo.GetActorsByMovie error while getting from movies_actors")
	}

	var actors []models.Actor

	tx = mr.DB.Table("actors").Find(&actors, actorIDs)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgMovieRepo.GetActorsByMovie error while getting actors")
	}

	return actors, nil
}

func (mr *pgMovieRepo) GetMoviesByTitle(title string) ([]models.Movie, error) {
	var movies []models.Movie

	title = "%" + title + "%"

	tx := mr.DB.Where("title LIKE ?", title).Find(&movies)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgMovieRepo.GetMoviesByTitle error while getting from movies")
	}

	return movies, nil
}
