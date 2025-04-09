package http

import (
	"log/slog"
	"net/http"

	pkgLogger "github.com/kostenbl4/code-tasks/pkg/log"
	"github.com/kostenbl4/code-tasks/task-service/internal/api/http/types"
	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/internal/usecases"
	"github.com/kostenbl4/code-tasks/task-service/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type UserHandler struct {
	logger *slog.Logger

	service  usecases.User
	smanager usecases.Session
}

func NewUserHandler(logger *slog.Logger, service usecases.User, smanager usecases.Session) *UserHandler {
	return &UserHandler{logger: logger, service: service, smanager: smanager}
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
func (uh *UserHandler) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	in, err := types.GetRegisterUserRequest(r)
	if err != nil {
		types.HandleError(w, domain.ErrBadRequest)
		return
	}

	log := uh.logger.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	_, err = uh.service.RegisterUser(in.Username, in.Password)
	if err != nil {
		log.Error("error while registering user: ", pkgLogger.Error(err))
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
func (uh *UserHandler) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	in, err := types.GetLoginUserRequest(r)
	if err != nil {
		types.HandleError(w, domain.ErrBadRequest)
		return
	}

	log := uh.logger.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	id, err := uh.service.LoginUser(in.Username, in.Password)

	if err != nil {
		log.Error("error while logging user in: ", pkgLogger.Error(err))
		types.HandleError(w, err)
		return
	}

	token, err := uh.smanager.CreateSession(id)

	if err != nil {
		log.Error("error while logging user in: ", pkgLogger.Error(err))
		types.HandleError(w, err)
		return
	}

	utils.WriteJSON(w, types.LoginUserResponse{Token: token}, http.StatusOK)
}

func (usrh *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/register", usrh.registerUserHandler)
	r.Post("/login", usrh.loginUserHandler)
}
