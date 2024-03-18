package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"intern/internal/actor/repository"
	"intern/models"
	"intern/pkg/logger"
)

type pgActorRepo struct {
	Logger logger.Logger
	DB     *gorm.DB
}

func New(logger logger.Logger, db *gorm.DB) repository.ActorRepositoryI {
	return &pgActorRepo{
		Logger: logger,
		DB:     db,
	}
}

func (ar *pgActorRepo) Create(a *models.Actor) error {
	tx := ar.DB.Create(a)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgActorRepo.Create error while inserting in repo")
	}

	return nil
}

func (ar *pgActorRepo) Get(id int) (*models.Actor, error) {
	var a models.Actor
	tx := ar.DB.Where("id = ?", id).Take(&a)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgActorRepo.Get error")
	}

	return &a, nil
}

func (ar *pgActorRepo) Update(a *models.Actor) error {
	tx := ar.DB.Omit("id").Updates(a)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgActorRepo.Update error while inserting in repo")
	}

	return nil
}

func (ar *pgActorRepo) Delete(id int) error {
	tx := ar.DB.Delete(&models.Actor{}, id)

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "pgActorRepo.Delete error")
	}

	return nil
}

func (ar *pgActorRepo) GetMoviesByActor(id int) ([]models.Movie, error) {
	var movieIDs []int

	tx := ar.DB.Table("movies_actors").Select("movie_id").Where("actor_id = ?", id).Find(&movieIDs)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgActorRepo.GetMoviesByActor error while getting from movies_actors")
	}

	var movies []models.Movie

	tx = ar.DB.Table("movies").Find(&movies, movieIDs)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgActorRepo.GetMoviesByActor error while getting movies")
	}

	return movies, nil
}
