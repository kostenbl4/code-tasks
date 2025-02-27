package http

import (
	"math/rand"
	"net/http"
	"task-server/internal/api/http/types"
	"task-server/internal/middleware/auth"
	"task-server/internal/repository"
	"task-server/internal/usecases"
	"task-server/utils"
	"time"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	service  usecases.Task
	smanager usecases.Session
}

func NewTaskHandler(service usecases.Task, smanager usecases.Session) *TaskHandler {
	return &TaskHandler{service: service, smanager: smanager}
}

// someHeavyWork - возвращает рандомный набор байт через продолжительное время
// - имитация работы другого сервиса с вычислениями
func someHeavyWork() []byte {
	time.Sleep(time.Second * time.Duration(3+rand.Intn(10)))
	res := make([]byte, 100)
	for i := 0; i < len(res); i++ {
		res[i] = byte(rand.Intn(100))
	}
	return res
}

// @Summary Create a task
// @Description Creates a new task and returns task ID
// @Tags task
// @Produce json
// @Success 201 {object} types.CreateTaskResponse
// @Failure 400 {object} types.ErrorResponse "Invalid request"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /task [post]
func (th *TaskHandler) createTaskHandler(w http.ResponseWriter, r *http.Request) {

	// var in types.CreateTaskRequest // читаем json в случае если передаются данные для обработки
	// err := utils.ReadJSON(r, &in)

	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, "failed to read json"+err.Error())
	// 	return
	// }

	u := th.service.CreateTask()

	go func() {
		res := someHeavyWork()
		t, err := th.service.GetTask(u)
		if err != nil {
			return
		}
		t.Status = "ready"
		t.Result = res
		err = th.service.UpdateTask(t)
		if err != nil {
			return
		}
	}()

	utils.WriteJSON(w, types.CreateTaskResponse{UUID: u.String()}, http.StatusCreated)
}

// @Summary Get task status
// @Description Returns the current status of a task by its ID
// @Tags task
// @Produce json
// @Param task_id path string true "Task ID" format(uuid)
// @Success 200 {object} types.GetTaskStatusResponse
// @Failure 400 {object} types.ErrorResponse "Invalid task ID format"
// @Failure 404 {object} types.ErrorResponse "Task not found"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /status/{task_id} [get]
func (th *TaskHandler) getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUID(r, "task_id")
	if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Invalid task_id"}, http.StatusBadRequest)
		return
	}

	task, err := th.service.GetTask(id)
	if err == repository.ErrTaskNotFound {
		utils.WriteJSON(w, types.ErrorResponse{Error: err.Error()}, http.StatusNotFound)
		return
	} else if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Internal error"}, http.StatusInternalServerError)
		return
	}
	status := task.Status

	utils.WriteJSON(w, types.GetTaskStatusResponse{Status: status}, http.StatusOK)
}

// @Summary Get task result
// @Description Returns the result of a completed task by its ID
// @Tags task
// @Produce json
// @Param task_id path string true "Task ID" format(uuid)
// @Success 200 {object} types.GetTaskResultResponse
// @Failure 400 {object} types.ErrorResponse "Invalid task ID format or internal error"
// @Failure 404 {object} types.ErrorResponse "Task not found"
// @Router /result/{task_id} [get]
func (th *TaskHandler) getTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUID(r, "task_id")
	if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Invalid task_id"}, http.StatusBadRequest)
		return
	}

	task, err := th.service.GetTask(id)
	if err == repository.ErrTaskNotFound {
		utils.WriteJSON(w, types.ErrorResponse{Error: err.Error()}, http.StatusNotFound)
		return
	} else if err != nil {
		utils.WriteJSON(w, types.ErrorResponse{Error: "Internal error"}, http.StatusInternalServerError)
		return
	}
	result := task.Result

	utils.WriteJSON(w, types.GetTaskResultResponse{Result: result}, http.StatusOK)
}

// RegisterRoutes - регистрация ручек
func (th *TaskHandler) RegisterRoutes(r chi.Router) {

	r.Group(func(r chi.Router) {
		r.Use(auth.SessionMiddleware(th.smanager))

		r.Post("/task", th.createTaskHandler)
		r.Get("/status/{task_id}", th.getTaskStatusHandler)
		r.Get("/result/{task_id}", th.getTaskResultHandler)
	})
}
