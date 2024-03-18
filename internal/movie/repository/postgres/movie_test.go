package postgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	movieRep "intern/internal/movie/repository"
	"intern/internal/testBuilders"
	"intern/models"
	"intern/pkg/logger"
	"regexp"
	"testing"
	"time"
)

type MovieRepoTestSuite struct {
	suite.Suite
	db           *sql.DB
	gormDB       *gorm.DB
	mock         sqlmock.Sqlmock
	repo         movieRep.MovieRepositoryI
	movieBuilder *testBuilders.MovieBuilder
}

func TestMovieRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(MovieRepoTestSuite))
}

func (s *MovieRepoTestSuite) BeforeEach(t provider.T) {
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
	s.movieBuilder = testBuilders.NewMovieBuilder()
}

func (s *MovieRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *MovieRepoTestSuite) TestCreateMovie(t provider.T) {
	movie := s.movieBuilder.
		WithID(1).
		WithTitle("movie").
		WithDesc("desc").
		WithRating(1).
		WithRelease(time.Date(1934, 1, 1, 0, 0, 0, 0, time.UTC)).
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "movies" ("title","description","release_date","rating","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(movie.Title, movie.Description, movie.ReleaseDate, movie.Rating, movie.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectCommit()

	err := s.repo.Create(&movie)
	t.Assert().NoError(err)
	t.Assert().Equal(1, movie.ID)
}

func (s *MovieRepoTestSuite) TestGetMovie(t provider.T) {
	movie := s.movieBuilder.
		WithID(1).
		WithTitle("movie").
		WithDesc("desc").
		WithRating(1).
		WithRelease(time.Date(1934, 1, 1, 0, 0, 0, 0, time.UTC)).
		Build()

	rows := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "rating"}).
		AddRow(
			movie.ID,
			movie.Title,
			movie.Description,
			movie.ReleaseDate,
			movie.Rating,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "movies" WHERE id = $1 LIMIT $2`)).
		WithArgs(movie.ID, 1).
		WillReturnRows(rows)

	resMovie, err := s.repo.Get(movie.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(movie, *resMovie)
}

func (s *MovieRepoTestSuite) TestUpdateMovie(t provider.T) {
	movie := s.movieBuilder.
		WithID(1).
		WithTitle("movie").
		WithDesc("desc").
		WithRating(1).
		WithRelease(time.Date(1934, 1, 1, 0, 0, 0, 0, time.UTC)).
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "movies" SET "title"=$1,"description"=$2,"release_date"=$3,"rating"=$4 WHERE "id" = $5`)).
		WithArgs(movie.Title, movie.Description, movie.ReleaseDate, movie.Rating, movie.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	err := s.repo.Update(&movie)
	t.Assert().NoError(err)
}

func (s *MovieRepoTestSuite) TestDeleteMovie(t provider.T) {
	movie := s.movieBuilder.
		WithID(1).
		WithTitle("movie").
		WithDesc("desc").
		WithRating(1).
		WithRelease(time.Date(1934, 1, 1, 0, 0, 0, 0, time.UTC)).
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "movies" WHERE "movies"."id" = $1`)).
		WithArgs(movie.ID).WillReturnResult(sqlmock.NewResult(int64(movie.ID), 1))

	s.mock.ExpectCommit()

	err := s.repo.Delete(movie.ID)
	t.Assert().NoError(err)
}

func (s *MovieRepoTestSuite) TestGetActorsByMovie(t provider.T) {
	actors := make([]models.Actor, 10)
	err := faker.FakeData(&actors)
	t.Assert().NoError(err)

	movieActors := make([]models.MovieActor, 10)
	err = faker.FakeData(&movieActors)
	t.Assert().NoError(err)

	for i := range movieActors {
		movieActors[i].MovieID = movieActors[0].MovieID
	}

	rows := sqlmock.NewRows([]string{"movie_id"})

	for i := range movieActors {
		rows.AddRow(movieActors[i].MovieID)
	}

	rowsActors := sqlmock.NewRows([]string{"id", "first_name", "last_name", "gender", "birthday"})

	for i := range actors {
		rowsActors.AddRow(movieActors[i].ID, actors[i].FirstName, actors[i].LastName, actors[i].Gender, actors[i].Birthday)
	}

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`select actor_id from movies_actors where movie_id = $1`)).
		WithArgs(movieActors[0].MovieID).
		WillReturnRows(rows)

	s.mock.ExpectCommit()

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`select * from actors where id IN $1`)).
		WithArgs(rows).
		WillReturnRows(rowsActors)

	resActors, err := s.repo.GetActorsByMovie(movieActors[0].MovieID)
	t.Assert().NoError(err)
	t.Assert().Equal(actors, resActors)
}
