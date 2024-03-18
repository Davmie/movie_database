package delivery

import (
	"encoding/json"
	movieUseCase "intern/internal/movie/usecase"
	"io"
	"net/http"
	"strconv"

	"intern/models"
	"intern/pkg/logger"

	"github.com/asaskevich/govalidator"
)

type MovieHandler struct {
	MovieUseCase movieUseCase.MovieUseCaseI
	Logger       logger.Logger
}

// Create godoc
// @Summary      Create a movie
// @Description  Create a movie
// @Tags     movies
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param    movie body models.Movie true "movie info"
// @Success 201 {object} models.Movie "movie created"
// @Failure 400 {object} http.Error "invalid body"
// @Failure 401 {object} http.Error "no auth"
// @Failure 403 {object} http.Error "forbidden"
// @Failure 500 {object} http.Error "internal server error"
// @Router   /movies/create [post]
func (mh *MovieHandler) Create(w http.ResponseWriter, r *http.Request) {
	movie := models.Movie{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		mh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		mh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &movie)
	if err != nil {
		mh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(movie)
	if err != nil {
		mh.Logger.Infow("can`t validate form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = mh.MovieUseCase.Create(&movie)
	if err != nil {
		mh.Logger.Infow("can`t create movie",
			"err:", err.Error())
		http.Error(w, "can`t create movie", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(movie)

	if err != nil {
		mh.Logger.Errorw("can`t marshal movie",
			"err:", err.Error())
		http.Error(w, "can`t make movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(resp)
	if err != nil {
		mh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// Get godoc
// @Summary      Get movie
// @Description  Get info about a movie by id
// @Tags     movies
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param id path int true "MOV_ID"
// @Success 200 {object} models.Movie "success get movie"
// @Failure 401 {object} http.Error "no auth"
// @Failure 403 {object} http.Error "forbidden"
// @Failure 404 {object} http.Error "Movie not found"
// @Failure 500 {object} http.Error "internal server error"
// @Router   /movies/{id} [get]
func (mh *MovieHandler) Get(w http.ResponseWriter, r *http.Request) {
	movieIdString := r.PathValue("MOV_ID")
	if movieIdString == "" {
		mh.Logger.Errorw("no MOV_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	movieId, err := strconv.Atoi(movieIdString)
	if err != nil {
		mh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	movie, err := mh.MovieUseCase.Get(movieId)
	if err != nil {
		mh.Logger.Infow("can`t get movie",
			"err:", err.Error())
		http.Error(w, "can`t get movie", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(movie)

	if err != nil {
		mh.Logger.Errorw("can`t marshal movie",
			"err:", err.Error())
		http.Error(w, "can`t make movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		mh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// Update godoc
// @Summary      Update movie
// @Description  Update info about a movie by id
// @Tags     movies
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param id path int true "MOV_ID"
// @Param movie body models.Movie true "movie info"
// @Success 200 {object} models.Movie "Movie updated"
// @Failure 400 {object} http.Error "invalid body"
// @Failure 401 {object} http.Error "no auth"
// @Failure 403 {object} http.Error "forbidden"
// @Failure 404 {object} http.Error "Movie not found"
// @Failure 500 {object} http.Error "internal server error"
// @Router   /movies/{id} [put]
func (mh *MovieHandler) Update(w http.ResponseWriter, r *http.Request) {
	movieIdString := r.PathValue("MOV_ID")
	if movieIdString == "" {
		mh.Logger.Errorw("no MOV_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	movieId, err := strconv.Atoi(movieIdString)
	if err != nil {
		mh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	movie := &models.Movie{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		mh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		mh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, movie)
	if err != nil {
		mh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(movie)
	if err != nil {
		mh.Logger.Infow("can`t validate form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	movie.ID = movieId
	err = mh.MovieUseCase.Update(movie)
	if err != nil {
		mh.Logger.Infow("can`t update movie",
			"err:", err.Error())
		http.Error(w, "can`t update movie", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(movie)

	if err != nil {
		mh.Logger.Errorw("can`t marshal movie",
			"err:", err.Error())
		http.Error(w, "can`t make movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		mh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// Delete godoc
// @Summary      Delete movie
// @Description  Delete info about a movie by id
// @Tags     movies
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param id path int true "MOV_ID"
// @Success 200 {object} models.Movie "Movie deleted"
// @Failure 401 {object} http.Error "no auth"
// @Failure 403 {object} http.Error "forbidden"
// @Failure 404 {object} http.Error "Movie not found"
// @Failure 500 {object} http.Error "internal server error"
// @Router   /movies/{id} [delete]
func (mh *MovieHandler) Delete(w http.ResponseWriter, r *http.Request) {
	movieIdString := r.PathValue("MOV_ID")
	if movieIdString == "" {
		mh.Logger.Errorw("no MOV_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	movieId, err := strconv.Atoi(movieIdString)
	if err != nil {
		mh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = mh.MovieUseCase.Delete(movieId)
	if err != nil {
		mh.Logger.Infow("can`t delete movie",
			"err:", err.Error())
		http.Error(w, "can`t delete movie", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetMoviesSorted godoc
// @Summary      Get sorted movies
// @Description  Get list of movies sorted by specified column
// @Tags     movies
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param id path int true "MOV_ID"
// @Success 200 {object} []models.Actor "success get sorted movies"
// @Failure 401 {object} http.Error "no auth"
// @Failure 403 {object} http.Error "forbidden"
// @Failure 404 {object} http.Error "Movie not found"
// @Failure 500 {object} http.Error "internal server error"
// @Router   /movies/sorted [get]
func (mh *MovieHandler) GetMoviesSorted(w http.ResponseWriter, r *http.Request) {
	column := r.FormValue("column")
	if column == "" {
		mh.Logger.Errorw("no column key")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	movies, err := mh.MovieUseCase.GetMoviesSorted(column)
	if err != nil {
		mh.Logger.Infow("can`t get movies",
			"err:", err.Error())
		http.Error(w, "can`t get movies", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(movies)

	if err != nil {
		mh.Logger.Errorw("can`t marshal movies",
			"err:", err.Error())
		http.Error(w, "can`t make movies", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		mh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// GetActorsByMovie godoc
// @Summary      Get movies' actors
// @Description  Get list of movies' actors by id
// @Tags     movies
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param id path int true "MOV_ID"
// @Success 200 {object} []models.Actor "success get actors by movie"
// @Failure 401 {object} http.Error "no auth"
// @Failure 403 {object} http.Error "forbidden"
// @Failure 404 {object} http.Error "Movie not found"
// @Failure 500 {object} http.Error "internal server error"
// @Router   /movies/{id}/actors [get]
func (mh *MovieHandler) GetActorsByMovie(w http.ResponseWriter, r *http.Request) {
	movieIdString := r.PathValue("MOV_ID")
	if movieIdString == "" {
		mh.Logger.Errorw("no MOV_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	movieId, err := strconv.Atoi(movieIdString)
	if err != nil {
		mh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	actors, err := mh.MovieUseCase.GetActorsByMovie(movieId)
	if err != nil {
		mh.Logger.Infow("can`t get actors",
			"err:", err.Error())
		http.Error(w, "can`t get actors", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(actors)

	if err != nil {
		mh.Logger.Errorw("can`t marshal actors",
			"err:", err.Error())
		http.Error(w, "can`t make actors", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		mh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// GetMoviesByTitle godoc
// @Summary      Get movies by title
// @Description  Get list of movies by fragment of title
// @Tags     movies
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param title query int true "title"
// @Success 200 {object} []models.Actor "success get movies by title"
// @Failure 401 {object} http.Error "no auth"
// @Failure 403 {object} http.Error "forbidden"
// @Failure 404 {object} http.Error "Movies not found"
// @Failure 500 {object} http.Error "internal server error"
// @Router   /movies/title [get]
func (mh *MovieHandler) GetMoviesByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if title == "" {
		mh.Logger.Errorw("no title key")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	movies, err := mh.MovieUseCase.GetMoviesByTitle(title)
	if err != nil {
		mh.Logger.Infow("can`t get movies",
			"err:", err.Error())
		http.Error(w, "can`t get movies", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(movies)

	if err != nil {
		mh.Logger.Errorw("can`t marshal movies",
			"err:", err.Error())
		http.Error(w, "can`t make movies", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		mh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}
