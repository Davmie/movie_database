package main

import (
	"fmt"
	"intern/cmd/server"
	actorDel "intern/internal/actor/delivery"
	pgActor "intern/internal/actor/repository/postgres"
	actorUseCase "intern/internal/actor/usecase"
	movieDel "intern/internal/movie/delivery"
	pgMovie "intern/internal/movie/repository/postgres"
	movieUseCase "intern/internal/movie/usecase"
	userDel "intern/internal/user/delivery"
	pgUser "intern/internal/user/repository/postgres"
	userUseCase "intern/internal/user/usecase"
	"intern/pkg/context"
	"intern/pkg/middleware"
	"intern/pkg/session"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var prodCfgPg = postgres.Config{DSN: "host=db user=postgres password=postgres port=5432"}

// @title MovieDataBase Swagger API
// @version 1.0
// @host localhost:8085
func main() {
	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	db, err := gorm.Open(postgres.New(prodCfgPg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sessionManager := session.JWTSessionsManager{}
	contextManager := context.Manager{}

	authManager := middleware.AuthManager{
		SessionManager: sessionManager,
		Logger:         logger,
		ContextManager: contextManager,
	}

	actorHandler := actorDel.ActorHandler{
		ActorUseCase: actorUseCase.New(pgActor.New(logger, db)),
		Logger:       logger,
	}

	movieHandler := movieDel.MovieHandler{
		MovieUseCase: movieUseCase.New(pgMovie.New(logger, db)),
		Logger:       logger,
	}

	userHandler := userDel.UserHandler{
		UserUseCase: userUseCase.New(pgUser.New(logger, db)),
		Logger:      logger,
		Sessions:    sessionManager,
	}

	r := http.NewServeMux()

	r.HandleFunc("POST /users/login", userHandler.Login)

	r.Handle("GET /actors/{ACT_ID}", authManager.Auth(http.HandlerFunc(actorHandler.Get), "user", "admin"))
	r.Handle("POST /actors", authManager.Auth(http.HandlerFunc(actorHandler.Create), "admin"))
	r.Handle("PUT /actors/{ACT_ID}", authManager.Auth(http.HandlerFunc(actorHandler.Update), "admin"))
	r.Handle("DELETE /actors/{ACT_ID}", authManager.Auth(http.HandlerFunc(actorHandler.Delete), "admin"))
	r.Handle("GET /actors/{ACT_ID}/movies", authManager.Auth(http.HandlerFunc(actorHandler.GetMoviesByActor), "user", "admin"))

	r.Handle("GET /movies/{MOV_ID}", authManager.Auth(http.HandlerFunc(movieHandler.Get), "user", "admin"))
	r.Handle("POST /movies", authManager.Auth(http.HandlerFunc(movieHandler.Create), "admin"))
	r.Handle("PUT /movies/{MOV_ID}", authManager.Auth(http.HandlerFunc(movieHandler.Update), "admin"))
	r.Handle("DELETE /movies/{MOV_ID}", authManager.Auth(http.HandlerFunc(movieHandler.Delete), "admin"))
	r.Handle("GET /movies/{MOV_ID}/actors", authManager.Auth(http.HandlerFunc(movieHandler.GetActorsByMovie), "user", "admin"))
	r.Handle("GET /movies/sorted", authManager.Auth(http.HandlerFunc(movieHandler.GetMoviesSorted), "user", "admin"))
	r.Handle("GET /movies/title", authManager.Auth(http.HandlerFunc(movieHandler.GetMoviesByTitle), "user", "admin"))

	router := middleware.AccessLog(logger, r)
	router = middleware.Panic(logger, router)

	s := server.NewServer(router)
	if err := s.Start(); err != nil {
		logger.Fatal(err)
	}

	err = zapLogger.Sync()
	if err != nil {
		fmt.Println(err)
	}
}
