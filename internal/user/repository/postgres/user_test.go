package postgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"intern/internal/testBuilders"
	userRep "intern/internal/user/repository"
	"intern/pkg/logger"
	"regexp"
	"testing"
)

type UserRepoTestSuite struct {
	suite.Suite
	db          *sql.DB
	gormDB      *gorm.DB
	mock        sqlmock.Sqlmock
	repo        userRep.UserRepositoryI
	userBuilder *testBuilders.UserBuilder
}

func TestUserRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(UserRepoTestSuite))
}

func (s *UserRepoTestSuite) BeforeEach(t provider.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("error while creating sql mock")
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatal("error gorm open")
	}

	var logger logger.Logger

	s.db = db
	s.gormDB = gormDB
	s.mock = mock

	s.repo = New(logger, gormDB)
	s.userBuilder = testBuilders.NewUserBuilder()
}

func (s *UserRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *UserRepoTestSuite) TestGetUser(t provider.T) {
	user := s.userBuilder.
		WithLogin("log").
		WithPassword("qsc").
		WithRole("user").
		Build()

	rows := sqlmock.NewRows([]string{"id", "name", "login", "password", "role"}).
		AddRow(
			user.ID,
			user.Login,
			user.Password,
			user.Role,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`select * from users where login = $1 and password = $2`)).
		WithArgs(user.Login, user.Password).
		WillReturnRows(rows)

	resUser, err := s.repo.GetByLoginAndPassword(user.Login, user.Password)
	t.Assert().NoError(err)
	t.Assert().Equal(user, resUser)
}
