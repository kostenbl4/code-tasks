package api

import (
	"math/rand"
	"net/http"
	"task-server/internal/storage"
	"task-server/utils"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// Server - структура сервера
type Server struct {
	Config Config
	Store  storage.Storage
}

// Config - структура конфигурации сервера
type Config struct {
	Addr string
}

// Run - запуск сервера
func (s *Server) Run() error {
	r := s.registerRoutes()

	srv := http.Server{
		Addr:    s.Config.Addr,
		Handler: r,
	}
	return srv.ListenAndServe()
}

// registerRoutes - регистрация ручек
func (s *Server) registerRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/task", s.handleTaskCreate)

	r.Get("/status/{task_id}", s.handleGetTaskStatus)

	r.Get("/result/{task_id}", s.handleGetTaskResult)

	return r
}

// someInput - структура для входных данных
type someInput struct {
	data string
}

// someHeavyWork - возвращает рандомный набор байт через продолжительное время
// - имитация работы другого сервиса с вычислениями
func someHeavyWork(data string) []byte {
	_ = data
	time.Sleep(time.Second * time.Duration(3+rand.Intn(30)))
	res := make([]byte, 100)
	for i := 0; i < len(res); i++ {
		res[i] = byte(rand.Intn(100))
	}
	return res
}

// handleTaskCreate - обработчик создания задачи
func (s *Server) handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	var in someInput
	err := utils.ReadJSON(r, &in)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "failed to read json"+err.Error())
		return
	}
	u := s.Store.CreateTask()

	go func() {
		res := someHeavyWork(in.data)
		t, err := s.Store.GetTask(u)
		if err != nil {
			return
		}
		t.Status = "ready"
		t.Result = res
		err = s.Store.UpdateTask(t)
		if err != nil {
			return
		}
	}()

	utils.WriteJSON(w, map[string]string{"UUID": u.String()}, http.StatusCreated)
}

// handleGetTaskStatus - обработчик получения статуса задачи
func (s *Server) handleGetTaskStatus(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "task_id")
	uuid, err := uuid.Parse(u)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	task, err := s.Store.GetTask(uuid)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	status := task.Status
	utils.WriteJSON(w, map[string]string{"status": status}, http.StatusOK)
}

// handleGetTaskResult - обработчик получения результата задачи
func (s *Server) handleGetTaskResult(w http.ResponseWriter, r *http.Request) {
	u := chi.URLParam(r, "task_id")
	uuid, err := uuid.Parse(u)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}
	task, err := s.Store.GetTask(uuid)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	result := task.Result
	utils.WriteJSON(w, map[string][]byte{"result": result}, http.StatusOK)
}
