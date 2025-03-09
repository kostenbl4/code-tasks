package types

import "task-server/internal/domain"

// CreateTaskRequest - структура для входных данных на создание задачи
type CreateTaskRequest struct {
	Translator string `json:"translator"`
	Code       string `json:"code"`
}

// CreateTaskResponse - структура для выходных данных на создание задачи
type CreateTaskResponse struct {
	UUID string `json:"task_id"`
}

// GetTaskStatusResponse - структура для выходных данных на запрос статуса задачи
type GetTaskStatusResponse struct {
	Status string `json:"status"`
}

// GetTaskResultResponse - структура для выходных данных на запрос результата задачи
type GetTaskResultResponse struct {
	Result TaskResult `json:"result"`
}

type TaskResult struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func CreateGetTaskResultResponse(task domain.Task) GetTaskResultResponse {
	return GetTaskResultResponse{
		TaskResult{
			task.Stdout,
			task.Stderr,
		},
	}
}
