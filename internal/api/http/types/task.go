package types

// CreateTaskRequest - структура для входных данных на создание задачи
// type CreateTaskRequest struct {
// 	Data string `json:"data"`
// }

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
	Result []byte `json:"result"`
}

