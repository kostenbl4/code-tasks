package http

import (
	"net/http"
	"task-server/internal/api/http/types"
	"task-server/internal/repository"
	"task-server/internal/usecases"
	"task-server/utils"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service  usecases.User
	smanager usecases.Session
}

func NewUserHandler(service usecases.User, smanager usecases.Session) *UserHandler {
	return &UserHandler{service: service, smanager: smanager}
}

// @Summary Register a new user
// @Description Registers a new user with a username and password
// @Tags user
// @Accept json
// @Produce json
// @Param user body types.RegisterUserRequest true "User credentials"
// @Success 201 {string} string "User registered successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid JSON"
// @Failure 500 {object} types.ErrorResponse "Internal error or user already exists"
// @Router /register [post]
func (usrh *UserHandler) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var in types.RegisterUserRequest

	err := utils.ReadJSON(r, &in)
	if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Invalid json"}, http.StatusBadRequest)
		return
	}

	_, err = usrh.service.RegisterUser(in.Username, in.Password)
	if err == usecases.ErrUserAlreadyExists {
		utils.WriteJSON(w, types.ErrorResponse{Error: err.Error()}, http.StatusInternalServerError)
		return
	} else if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Internal error"}, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Login a user
// @Description Logs in a user and returns a session token
// @Tags user
// @Accept json
// @Produce json
// @Param user body types.LoginUserRequest true "User credentials"
// @Success 200 {object} types.LoginUserResponse "User logged in successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid JSON"
// @Failure 401 {object} types.ErrorResponse "Incorrect login or password"
// @Failure 303 {object} types.ErrorResponse "User already logged in"
// @Failure 500 {object} types.ErrorResponse "Internal error"
// @Router /login [post]
func (usrh *UserHandler) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var in types.LoginUserRequest

	err := utils.ReadJSON(r, &in)
	if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Invalid json"}, http.StatusBadRequest)
		return
	}

	id, err := usrh.service.LoginUser(in.Username, in.Password)
	if err == usecases.ErrIncorrectPassword || err == repository.ErrUserNotFound {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Incorrect login or password"}, http.StatusUnauthorized)
		return
	} else if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Internal error"}, http.StatusInternalServerError)
		return
	}

	token, err := usrh.smanager.CreateSession(id)
	if err == usecases.ErrSessionAlreadyExists {
		utils.WriteJSON(w, types.ErrorResponse{Error: "User already logged in"}, http.StatusSeeOther)
		return
	} else if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Internal error"}, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, types.LoginUserResponse{Token: token}, http.StatusOK)
}

func (usrh *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", usrh.registerUserHandler)
	r.Post("/login", usrh.loginUserHandler)
}
