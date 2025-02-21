package http

import (
	"math/rand"
	"net/http"
	"task-server/internal/api/http/types"
	"task-server/internal/repository"
	"task-server/internal/usecases"
	"task-server/utils"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type TaskHandler struct {
	service usecases.Task
}

func NewTaskHandler(service usecases.Task) *TaskHandler {
	return &TaskHandler{service: service}
}

// someHeavyWork - возвращает рандомный набор байт через продолжительное время
// - имитация работы другого сервиса с вычислениями
func someHeavyWork() []byte {
	time.Sleep(time.Second * time.Duration(3+rand.Intn(30)))
	res := make([]byte, 100)
	for i := 0; i < len(res); i++ {
		res[i] = byte(rand.Intn(100))
	}
	return res
}

// createTaskHandler - обработчик создания задачи
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

// getTaskStatusHandler - обработчик получения статуса задачи
func (th *TaskHandler) getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUID(r, "task_id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid task_id")
		return
	}

	task, err := th.service.GetTask(id)
	if err == repository.ErrTaskNotFound {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Internal error")
		return
	}
	status := task.Status

	utils.WriteJSON(w, types.GetStatusResponse{Status: status}, http.StatusOK)
}

// getTaskResultHandler - обработчик получения результата задачи
func (th *TaskHandler) getTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseUUID(r, "task_id")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid task_id")
		return
	}

	task, err := th.service.GetTask(id)
	if err == repository.ErrTaskNotFound {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	} else if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Internal error")
		return
	}
	result := task.Result

	utils.WriteJSON(w, types.GetResultResponse{Result: result}, http.StatusOK)
}

// RegisterRoutes - регистрация ручек
func (th *TaskHandler) RegisterRoutes(r chi.Router) {

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/task", th.createTaskHandler)

	r.Get("/status/{task_id}", th.getTaskStatusHandler)

	r.Get("/result/{task_id}", th.getTaskResultHandler)
}
