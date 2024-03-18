package postgres

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"intern/internal/user/repository"
	"intern/models"
	"intern/pkg/logger"
)

type pgUserRepo struct {
	Logger logger.Logger
	DB     *gorm.DB
}

func New(logger logger.Logger, db *gorm.DB) repository.UserRepositoryI {
	return &pgUserRepo{
		Logger: logger,
		DB:     db,
	}
}

func (ur *pgUserRepo) GetByLoginAndPassword(login, password string) (*models.User, error) {
	var u models.User
	tx := ur.DB.Where("login = ?", login).Where("password = ?", password).Take(&u)

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "pgUserRepo.GetByLoginAndPassword error")
	}

	return &u, nil
}
