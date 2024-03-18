package usecase

import (
	"github.com/pkg/errors"
	movieRep "intern/internal/movie/repository"
	"intern/models"
)

type MovieUseCaseI interface {
	Create(a *models.Movie) error
	Get(id int) (*models.Movie, error)
	Update(a *models.Movie) error
	Delete(id int) error
	GetMoviesSorted(sortingColumn string) ([]models.Movie, error)
	GetActorsByMovie(id int) ([]models.Actor, error)
	GetMoviesByTitle(title string) ([]models.Movie, error)
}

type movieUseCase struct {
	movieRepository movieRep.MovieRepositoryI
}

func New(aRep movieRep.MovieRepositoryI) MovieUseCaseI {
	return &movieUseCase{
		movieRepository: aRep,
	}
}

func (mUC *movieUseCase) Create(a *models.Movie) error {
	err := mUC.movieRepository.Create(a)

	if err != nil {
		return errors.Wrap(err, "movieUseCase.Create error")
	}

	return nil
}

func (mUC *movieUseCase) Get(id int) (*models.Movie, error) {
	resMovie, err := mUC.movieRepository.Get(id)

	if err != nil {
		return nil, errors.Wrap(err, "movieUseCase.Get error")
	}

	return resMovie, nil
}

func (mUC *movieUseCase) Update(a *models.Movie) error {
	_, err := mUC.movieRepository.Get(a.ID)

	if err != nil {
		return errors.Wrap(err, "movieUseCase.Update error: Movie not found")
	}

	err = mUC.movieRepository.Update(a)

	if err != nil {
		return errors.Wrap(err, "movieUseCase.Update error: Can't update in repo")
	}

	return nil
}

func (mUC *movieUseCase) Delete(id int) error {
	_, err := mUC.movieRepository.Get(id)

	if err != nil {
		return errors.Wrap(err, "movieUseCase.Delete error: Movie not found")
	}

	err = mUC.movieRepository.Delete(id)

	if err != nil {
		return errors.Wrap(err, "movieUseCase.Delete error: Can't delete in repo")
	}

	return nil
}

func (mUC *movieUseCase) GetMoviesSorted(sortingColumn string) ([]models.Movie, error) {
	movies, err := mUC.movieRepository.GetMoviesSorted(sortingColumn)

	if err != nil {
		return nil, errors.Wrap(err, "movieUseCase.GetMoviesSorted error")
	}

	return movies, nil
}

func (mUC *movieUseCase) GetActorsByMovie(id int) ([]models.Actor, error) {
	actors, err := mUC.movieRepository.GetActorsByMovie(id)

	if err != nil {
		return nil, errors.Wrap(err, "movieUseCase.GetMoviesByMovie error")
	}

	return actors, nil
}

func (mUC *movieUseCase) GetMoviesByTitle(title string) ([]models.Movie, error) {
	movies, err := mUC.movieRepository.GetMoviesByTitle(title)

	if err != nil {
		return nil, errors.Wrap(err, "movieUseCase.GetMoviesByTitle error")
	}

	return movies, nil
}
