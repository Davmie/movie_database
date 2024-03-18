package delivery

import (
	"encoding/json"
	actorUseCase "intern/internal/actor/usecase"
	"io"
	"net/http"
	"strconv"

	"intern/models"
	"intern/pkg/logger"

	"github.com/asaskevich/govalidator"
)

type ActorHandler struct {
	ActorUseCase actorUseCase.ActorUseCaseI
	Logger       logger.Logger
}

// Create godoc
// @Summary      Create an actor
// @Description  Create an actor
// @Tags     actors
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param    actor body models.Actor true "actor info"
// @Success 201 {object} models.Actor "actor created"
// @Failure 400 {object} nil "invalid body"
// @Failure 401 {object} nil "no auth"
// @Failure 403 {object} nil "forbidden"
// @Failure 500 {object} nil "internal server error"
// @Router   /actors [post]
func (ah *ActorHandler) Create(w http.ResponseWriter, r *http.Request) {
	actor := models.Actor{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ah.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &actor)
	if err != nil {
		ah.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(actor)
	if err != nil {
		ah.Logger.Infow("can`t validate form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = ah.ActorUseCase.Create(&actor)
	if err != nil {
		ah.Logger.Infow("can`t create actor",
			"err:", err.Error())
		http.Error(w, "can`t create actor", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(actor)

	if err != nil {
		ah.Logger.Errorw("can`t marshal actor",
			"err:", err.Error())
		http.Error(w, "can`t make actor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// Get godoc
// @Summary      Get actor
// @Description  Get info about an actor by id
// @Tags     actors
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param id path int true "ACT_ID"
// @Success 200 {object} models.Actor "success get actor"
// @Failure 401 {object} nil "no auth"
// @Failure 403 {object} nil "forbidden"
// @Failure 404 {object} nil "Actor not found"
// @Failure 500 {object} nil "internal server error"
// @Router   /actors/{id} [get]
func (ah *ActorHandler) Get(w http.ResponseWriter, r *http.Request) {
	actorIdString := r.PathValue("ACT_ID")
	if actorIdString == "" {
		ah.Logger.Errorw("no ACT_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	actorId, err := strconv.Atoi(actorIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	actor, err := ah.ActorUseCase.Get(actorId)
	if err != nil {
		ah.Logger.Infow("can`t get actor",
			"err:", err.Error())
		http.Error(w, "can`t get actor", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(actor)

	if err != nil {
		ah.Logger.Errorw("can`t marshal actor",
			"err:", err.Error())
		http.Error(w, "can`t make actor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// Update godoc
// @Summary      Update actor
// @Description  Update info about an actor by id
// @Tags     actors
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param id path int true "ACT_ID"
// @Param actor body models.Actor true "actor info"
// @Success 200 {object} models.Actor "Actor updated"
// @Failure 400 {object} nil "invalid body"
// @Failure 401 {object} nil "no auth"
// @Failure 403 {object} nil "forbidden"
// @Failure 404 {object} nil "Actor not found"
// @Failure 500 {object} nil "internal server error"
// @Router   /actors/{id} [put]
func (ah *ActorHandler) Update(w http.ResponseWriter, r *http.Request) {
	actorIdString := r.PathValue("ACT_ID")
	if actorIdString == "" {
		ah.Logger.Errorw("no ACT_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	actorId, err := strconv.Atoi(actorIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	actor := &models.Actor{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		ah.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, actor)
	if err != nil {
		ah.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(actor)
	if err != nil {
		ah.Logger.Infow("can`t validate form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	actor.ID = actorId
	err = ah.ActorUseCase.Update(actor)
	if err != nil {
		ah.Logger.Infow("can`t update actor",
			"err:", err.Error())
		http.Error(w, "can`t update actor", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(actor)

	if err != nil {
		ah.Logger.Errorw("can`t marshal actor",
			"err:", err.Error())
		http.Error(w, "can`t make actor", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// Delete godoc
// @Summary      Delete actor
// @Description  Delete info about an actor by id
// @Tags     actors
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param id path int true "ACT_ID"
// @Success 200 {object} models.Actor "Actor deleted"
// @Failure 401 {object} nil "no auth"
// @Failure 403 {object} nil "forbidden"
// @Failure 404 {object} nil "Actor not found"
// @Failure 500 {object} nil "internal server error"
// @Router   /actors/{id} [delete]
func (ah *ActorHandler) Delete(w http.ResponseWriter, r *http.Request) {
	actorIdString := r.PathValue("ACT_ID")
	if actorIdString == "" {
		ah.Logger.Errorw("no ACT_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	actorId, err := strconv.Atoi(actorIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = ah.ActorUseCase.Delete(actorId)
	if err != nil {
		ah.Logger.Infow("can`t delete actor",
			"err:", err.Error())
		http.Error(w, "can`t delete actor", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetMoviesByActor godoc
// @Summary      Get actor's movies
// @Description  Get list of actor's movies by id
// @Tags     actors
// @Accept	 application/json
// @Produce  application/json
// @Param    Authorization header string true "token"
// @Param id path int true "ACT_ID"
// @Success 200 {object} []models.Actor "success get movies by actor"
// @Failure 401 {object} nil "no auth"
// @Failure 403 {object} nil "forbidden"
// @Failure 404 {object} nil "Actor not found"
// @Failure 500 {object} nil "internal server error"
// @Router   /actors/{id}/movies [get]
func (ah *ActorHandler) GetMoviesByActor(w http.ResponseWriter, r *http.Request) {
	actorIdString := r.PathValue("ACT_ID")
	if actorIdString == "" {
		ah.Logger.Errorw("no ACT_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	actorId, err := strconv.Atoi(actorIdString)
	if err != nil {
		ah.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	expCats, err := ah.ActorUseCase.GetMoviesByActor(actorId)
	if err != nil {
		ah.Logger.Infow("can`t get movies",
			"err:", err.Error())
		http.Error(w, "can`t get movies", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(expCats)

	if err != nil {
		ah.Logger.Errorw("can`t marshal movies",
			"err:", err.Error())
		http.Error(w, "can`t make movies", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		ah.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}
