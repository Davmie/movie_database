package repository

import "intern/models"

type UserRepositoryI interface {
	GetByLoginAndPassword(login, password string) (*models.User, error)
}
