package types

import "github.com/kostenbl4/code-tasks/task-service/internal/domain"

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
// пока что не знаю как учесть такую вложенность в swagger
type GetTaskResultResponse struct {
	Result string `json:"result"`
	Data   any    `json:"data"`
}

type TaskResultOK struct {
	Stdout string `json:"stdout"`
}

type TaskResponseError struct {
	Stderr string `json:"stderr"`
}

type CommitTaskRequest struct {
	domain.Task
}

func CreateGetTaskResultResponse(task domain.Task) GetTaskResultResponse {
	if task.Result == domain.TaskResultOk {
		return GetTaskResultResponse{
			Result: "ok",
			Data: TaskResultOK{
				Stdout: task.Stdout,
			},
		}
	} else {
		return GetTaskResultResponse{
			Result: task.Result,
			Data: TaskResponseError{
				Stderr: task.Stderr,
			},
		}
	}
}
