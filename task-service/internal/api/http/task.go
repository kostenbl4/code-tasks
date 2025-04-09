package http

import (
	"log/slog"
	"net/http"

	pkgLogger "github.com/kostenbl4/code-tasks/pkg/log"
	"github.com/kostenbl4/code-tasks/task-service/internal/api/http/types"
	"github.com/kostenbl4/code-tasks/task-service/internal/domain"
	"github.com/kostenbl4/code-tasks/task-service/internal/middleware/auth"
	"github.com/kostenbl4/code-tasks/task-service/internal/usecases"
	"github.com/kostenbl4/code-tasks/task-service/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type TaskHandler struct {
	logger *slog.Logger

	service  usecases.Task
	smanager usecases.Session
}

func NewTaskHandler(logger *slog.Logger, service usecases.Task, smanager usecases.Session) *TaskHandler {
	return &TaskHandler{logger: logger, service: service, smanager: smanager}
}

// @Summary Create a new task
// @Description Creates a new task with the provided translator and code, and returns the unique task ID.
// @Tags task
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}" default(Bearer <ваш_токен>)
// @Param task body types.CreateTaskRequest true "Task creation request payload"
// @Success 201 {object} types.CreateTaskResponse "Task successfully created"
// @Failure 400 {object} types.ErrorResponse "Invalid request payload"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /task [post]
func (th *TaskHandler) createTaskHandler(w http.ResponseWriter, r *http.Request) {

	var in types.CreateTaskRequest // читаем json в случае если передаются данные для обработки
	err := utils.ReadJSON(r, &in)
	if err != nil {
		types.HandleError(w, domain.ErrBadRequest)
		return
	}

	userID, err := utils.GetContextInt(r, auth.UserIDKey)
	if err != nil {
		types.HandleError(w, domain.ErrBadRequest)
		return
	}

	log := th.logger.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
		slog.Int("user_id", userID),
	)

	task, err := th.service.CreateTask(in.Translator, in.Code, int64(userID))
	if err != nil {
		log.Error("error while creating task: ", pkgLogger.Error(err))
		types.HandleError(w, err)
		return
	}

	err = th.service.SendTask(task)
	if err != nil {
		log.Error("error while sending task: ", pkgLogger.Error(err))
		types.HandleError(w, err)
		return
	}

	utils.WriteJSON(w, types.CreateTaskResponse{UUID: task.ID.String()}, http.StatusCreated)
}

// @Summary Retrieve task status
// @Description Fetches the current status of a task using its unique ID.
// @Tags task
// @Produce json
// @Param task_id path string true "Unique Task ID" format(uuid)
// @Param Authorization header string true "Bearer {token}" default(Bearer <ваш_токен>)
// @Success 200 {object} types.GetTaskStatusResponse "Current task status"
// @Failure 400 {object} types.ErrorResponse "Invalid task ID format"
// @Failure 404 {object} types.ErrorResponse "Task not found"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /status/{task_id} [get]
func (th *TaskHandler) getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUID(r, "task_id")
	if err != nil {
		types.HandleError(w, domain.ErrBadRequest)
		return
	}

	log := th.logger.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	task, err := th.service.GetTask(id)
	if err != nil {
		log.Error("error while getting task status: ", pkgLogger.Error(err))
		types.HandleError(w, err)
		return
	}

	utils.WriteJSON(w, types.GetTaskStatusResponse{Status: task.Status}, http.StatusOK)
}

// @Summary Retrieve task result
// @Description Returns the result of a completed task using its unique ID.
// @Tags task
// @Produce json
// @Param task_id path string true "Unique Task ID" format(uuid)
// @Param Authorization header string true "Bearer {token}" default(Bearer <ваш_токен>)
// @Success 200 {object} types.GetTaskResultResponse "Task result"
// @Failure 400 {object} types.ErrorResponse "Invalid task ID format or internal error"
// @Failure 404 {object} types.ErrorResponse "Task not found"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /result/{task_id} [get]
func (th *TaskHandler) getTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUID(r, "task_id")
	if err != nil {
		types.HandleError(w, domain.ErrBadRequest)
		return
	}

	log := th.logger.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	task, err := th.service.GetTask(id)
	if err != nil {
		log.Error("error while getting task result: ", pkgLogger.Error(err))
		types.HandleError(w, err)
		return
	}

	utils.WriteJSON(w, types.CreateGetTaskResultResponse(task), http.StatusOK)
}

// @Summary Commit task result
// @Description Commits the result of a task using its unique ID.
// @Tags task
// @Accept json
// @Produce json
// @Param task body types.CommitTaskRequest true "Task result commit request payload"
// @Success 200
// @Failure 400 {object} types.ErrorResponse "Invalid request payload"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /commit [put]
func (th *TaskHandler) commitTaskResult(w http.ResponseWriter, r *http.Request) {
	var in types.CommitTaskRequest
	err := utils.ReadJSON(r, &in)
	if err != nil {
		types.HandleError(w, domain.ErrBadRequest)
		return
	}

	log := th.logger.With(
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	task := domain.Task{
		ID:         in.ID,
		Status:     in.Status,
		Result:     in.Result,
		Translator: in.Translator,
		Code:       in.Code,
		Stdout:     in.Stdout,
		Stderr:     in.Stderr,
	}

	if err := th.service.UpdateTask(task); err != nil {
		log.Error("error while commiting task: ", pkgLogger.Error(err))
		types.HandleError(w, err)
		return
	}

}

// RegisterRoutes - регистрация ручек
func (th *TaskHandler) RegisterRoutes(r chi.Router) {

	r.Group(func(r chi.Router) {
		r.Use(auth.SessionMiddleware(th.smanager))

		r.Post("/task", th.createTaskHandler)
		r.Get("/status/{task_id}", th.getTaskStatusHandler)
		r.Get("/result/{task_id}", th.getTaskResultHandler)
	})

	// Для упрощения тестирования убрал аутентификацию с этого эндпоинта
	r.Put("/commit", th.commitTaskResult)
}
