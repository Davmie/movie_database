package delivery

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"io"
	"net/http"

	userUseCase "intern/internal/user/usecase"
	"intern/pkg/logger"
)

type TokenForm struct {
	Token string `json:"token"`
}

type LoginForm struct {
	Login    string `valid:"minstringlength(5)" json:"login"`
	Password string `valid:"minstringlength(5)" json:"password"`
}

type SessionManager interface {
	CreateSession(int, string) (string, error)
}

type UserHandler struct {
	UserUseCase userUseCase.UserUseCaseI
	Logger      logger.Logger
	Sessions    SessionManager
}

// Login godoc
// @Summary      User login
// @Description  User sign in with login and password
// @Tags     users
// @Accept	 application/json
// @Produce  application/json
// @Param    user body LoginForm true "user login and password"
// @Success 200 {object} TokenForm.Token "User signed in"
// @Failure 400 {object} http.Error "invalid body"
// @Failure 404 {object} http.Error "User not found"
// @Failure 500 {object} http.Error "internal server error"
// @Router   /users/login [post]
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	logForm := &LoginForm{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		uh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad reg data", http.StatusBadRequest)
		return
	}

	err = r.Body.Close()
	if err != nil {
		uh.Logger.Errorw("can`t close body of request", "err:", err.Error())
		http.Error(w, "close error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, logForm)
	if err != nil {
		uh.Logger.Infow("can`t unmarshal register form",
			"err:", err.Error())
		http.Error(w, "bad reg data", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(logForm)
	if err != nil {
		uh.Logger.Infow("can`t validate register form",
			"err:", err.Error())
		http.Error(w, "bad reg data", http.StatusBadRequest)
		return
	}

	user, err := uh.UserUseCase.GetByLoginAndPassword(logForm.Login, logForm.Password)
	if err != nil {
		uh.Logger.Infow("can`t get user by login and password",
			"err:", err.Error())
		http.Error(w, "can`t login", http.StatusNotFound)
		return
	}

	token, err := uh.Sessions.CreateSession(user.ID, user.Role)
	if err != nil {
		uh.Logger.Errorw("can`t create session",
			"err:", err.Error())
		http.Error(w, "can`t make session", http.StatusInternalServerError)
		return
	}

	tokenForm := &TokenForm{token}
	resp, err := json.Marshal(tokenForm)

	if err != nil {
		uh.Logger.Errorw("can`t marshal session token",
			"err:", err.Error())
		http.Error(w, "can`t make session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		uh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}
