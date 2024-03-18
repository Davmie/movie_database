package usecase

import (
	"github.com/pkg/errors"
	userRep "intern/internal/user/repository"
	"intern/models"
)

type UserUseCaseI interface {
	GetByLoginAndPassword(login, password string) (*models.User, error)
}

type userUseCase struct {
	userRepository userRep.UserRepositoryI
}

func New(uRep userRep.UserRepositoryI) UserUseCaseI {
	return &userUseCase{
		userRepository: uRep,
	}
}

func (u *userUseCase) GetByLoginAndPassword(login, password string) (*models.User, error) {
	resUser, err := u.userRepository.GetByLoginAndPassword(login, password)

	if err != nil {
		return nil, errors.Wrap(err, "userUseCase.GetByLoginAndPassword error")
	}

	return resUser, nil
}
