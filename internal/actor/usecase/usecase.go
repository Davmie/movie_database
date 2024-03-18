package usecase

import (
	"github.com/pkg/errors"
	actorRep "intern/internal/actor/repository"
	"intern/models"
)

type ActorUseCaseI interface {
	Create(a *models.Actor) error
	Get(id int) (*models.Actor, error)
	Update(a *models.Actor) error
	Delete(id int) error
	GetMoviesByActor(id int) ([]models.Movie, error)
}

type actorUseCase struct {
	actorRepository actorRep.ActorRepositoryI
}

func New(aRep actorRep.ActorRepositoryI) ActorUseCaseI {
	return &actorUseCase{
		actorRepository: aRep,
	}
}

func (aUC *actorUseCase) Create(a *models.Actor) error {
	err := aUC.actorRepository.Create(a)

	if err != nil {
		return errors.Wrap(err, "actorUseCase.Create error")
	}

	return nil
}

func (aUC *actorUseCase) Get(id int) (*models.Actor, error) {
	resActor, err := aUC.actorRepository.Get(id)

	if err != nil {
		return nil, errors.Wrap(err, "actorUseCase.Get error")
	}

	return resActor, nil
}

func (aUC *actorUseCase) Update(a *models.Actor) error {
	_, err := aUC.actorRepository.Get(a.ID)

	if err != nil {
		return errors.Wrap(err, "actorUseCase.Update error: Actor not found")
	}

	err = aUC.actorRepository.Update(a)

	if err != nil {
		return errors.Wrap(err, "actorUseCase.Update error: Can't update in repo")
	}

	return nil
}

func (aUC *actorUseCase) Delete(id int) error {
	_, err := aUC.actorRepository.Get(id)

	if err != nil {
		return errors.Wrap(err, "actorUseCase.Delete error: Actor not found")
	}

	err = aUC.actorRepository.Delete(id)

	if err != nil {
		return errors.Wrap(err, "actorUseCase.Delete error: Can't delete in repo")
	}

	return nil
}

func (aUC *actorUseCase) GetMoviesByActor(id int) ([]models.Movie, error) {
	movies, err := aUC.actorRepository.GetMoviesByActor(id)

	if err != nil {
		return nil, errors.Wrap(err, "actorUseCase.GetMoviesByActor error")
	}

	return movies, nil
}
