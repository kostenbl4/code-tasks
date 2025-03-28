package http

import (
	"code-tasks/task-service/internal/api/http/types"
	"code-tasks/task-service/internal/usecases"
	"code-tasks/task-service/utils"
	"net/http"

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
// @Failure 400 {object} types.ErrorResponse "Invalid JSON or user already exists"
// @Failure 500 {object} types.ErrorResponse "Internal error"
// @Router /register [post]
func (usrh *UserHandler) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	in, err := types.GetRegisterUserRequest(r)
	if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Invalid json"}, http.StatusBadRequest)
		return
	}

	_, err = usrh.service.RegisterUser(in.Username, in.Password)
	if err != nil {
		types.HandleError(w, err)
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
// @Failure 500 {object} types.ErrorResponse "Internal error"
// @Router /login [post]
func (usrh *UserHandler) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	in, err := types.GetLoginUserRequest(r)
	if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Invalid json"}, http.StatusBadRequest)
		return
	}

	id, err := usrh.service.LoginUser(in.Username, in.Password)
	
	if err != nil {
		types.HandleError(w, err)
		return
	}

	token, err := usrh.smanager.CreateSession(id)
	
	if err != nil {
		types.HandleError(w, err)
		return
	}

	utils.WriteJSON(w, types.LoginUserResponse{Token: token}, http.StatusOK)
}

func (usrh *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", usrh.registerUserHandler)
	r.Post("/login", usrh.loginUserHandler)
}
