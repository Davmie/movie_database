package postgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	actorRep "intern/internal/actor/repository"
	"intern/internal/testBuilders"
	"intern/models"
	"intern/pkg/logger"
	"regexp"
	"testing"
	"time"
)

type ActorRepoTestSuite struct {
	suite.Suite
	db           *sql.DB
	gormDB       *gorm.DB
	mock         sqlmock.Sqlmock
	repo         actorRep.ActorRepositoryI
	actorBuilder *testBuilders.ActorBuilder
}

func TestActorRepoSuite(t *testing.T) {
	suite.RunSuite(t, new(ActorRepoTestSuite))
}

func (s *ActorRepoTestSuite) BeforeEach(t provider.T) {
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
	s.actorBuilder = testBuilders.NewActorBuilder()
}

func (s *ActorRepoTestSuite) AfterEach(t provider.T) {
	err := s.mock.ExpectationsWereMet()
	t.Assert().NoError(err)
	s.db.Close()
}

func (s *ActorRepoTestSuite) TestCreateActor(t provider.T) {
	actor := s.actorBuilder.
		WithID(1).
		WithFirstName("act").
		WithLastName("act").
		WithGender('f').
		WithBirthday(time.Date(1934, 1, 1, 0, 0, 0, 0, time.UTC)).
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "actors" ("first_name","last_name","gender","birthday","id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(actor.FirstName, actor.LastName, actor.Gender, actor.Birthday, actor.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	s.mock.ExpectCommit()

	err := s.repo.Create(&actor)
	t.Assert().NoError(err)
	t.Assert().Equal(1, actor.ID)
}

func (s *ActorRepoTestSuite) TestGetActor(t provider.T) {
	birth := time.Date(1934, 1, 1, 0, 0, 0, 0, time.UTC)
	actor := s.actorBuilder.
		WithID(1).
		WithFirstName("act").
		WithLastName("act").
		WithGender('f').
		WithBirthday(birth).
		Build()

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "gender", "birthday"}).
		AddRow(
			actor.ID,
			actor.FirstName,
			actor.LastName,
			actor.Gender,
			actor.Birthday,
		)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "actors" WHERE id = $1 LIMIT $2`)).
		WithArgs(actor.ID, 1).
		WillReturnRows(rows)

	resActor, err := s.repo.Get(actor.ID)
	t.Assert().NoError(err)
	t.Assert().Equal(actor, *resActor)
}

func (s *ActorRepoTestSuite) TestUpdateActor(t provider.T) {
	birth := time.Date(1934, 1, 1, 0, 0, 0, 0, time.UTC)
	actor := s.actorBuilder.
		WithID(1).
		WithFirstName("act").
		WithLastName("act").
		WithGender('f').
		WithBirthday(birth).
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "actors" SET "first_name"=$1,"last_name"=$2,"gender"=$3,"birthday"=$4 WHERE "id" = $5`)).
		WithArgs(actor.FirstName, actor.LastName, actor.Gender, actor.Birthday, actor.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectCommit()

	err := s.repo.Update(&actor)
	t.Assert().NoError(err)
}

func (s *ActorRepoTestSuite) TestDeleteActor(t provider.T) {
	birth := time.Date(1934, 1, 1, 0, 0, 0, 0, time.UTC)
	actor := s.actorBuilder.
		WithID(1).
		WithFirstName("act").
		WithLastName("act").
		WithGender('f').
		WithBirthday(birth).
		Build()

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "actors" WHERE "actors"."id" = $1`)).
		WithArgs(actor.ID).WillReturnResult(sqlmock.NewResult(int64(actor.ID), 1))

	s.mock.ExpectCommit()

	err := s.repo.Delete(actor.ID)
	t.Assert().NoError(err)
}

func (s *ActorRepoTestSuite) TestGetMoviesByActor(t provider.T) {
	movies := make([]models.Movie, 10)
	err := faker.FakeData(&movies)
	t.Assert().NoError(err)

	movieActors := make([]models.MovieActor, 10)
	err = faker.FakeData(&movieActors)
	t.Assert().NoError(err)

	for i := range movieActors {
		movieActors[i].ActorID = movieActors[0].ActorID
	}

	rows := sqlmock.NewRows([]string{"movie_id"})

	for i := range movieActors {
		rows.AddRow(movieActors[i].MovieID)
	}

	rowsMovies := sqlmock.NewRows([]string{"id", "title", "description", "release_date", "rating"})

	for i := range movies {
		rowsMovies.AddRow(movieActors[i].ID, movies[i].Title, movies[i].Description, movies[i].ReleaseDate, movies[i].Rating)
	}

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`select movie_id from movies_actors where actor_id = $1`)).
		WithArgs(movieActors[0].ActorID).
		WillReturnRows(rows)

	s.mock.ExpectCommit()

	s.mock.ExpectBegin()

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`select * from movies where id IN $1`)).
		WithArgs(rows).
		WillReturnRows(rowsMovies)

	resMovies, err := s.repo.GetMoviesByActor(movieActors[0].ActorID)
	t.Assert().NoError(err)
	t.Assert().Equal(movies, resMovies)
}
