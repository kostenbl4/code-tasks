package types

// CreateTaskRequest - структура для входных данных на создание задачи
type CreateTaskRequest struct {
	Data string `json:"data"`
}

// CreateTaskResponse - структура для выходных данных на создание задачи
type CreateTaskResponse struct {
	UUID string `json:"uuid"`
}

// GetStatusResponse - структура для выходных данных на запрос статуса задачи
type GetStatusResponse struct {
	Status string `json:"status"`
}

// GetResultResponse - структура для выходных данных на запрос результата задачи
type GetResultResponse struct {
	Result []byte `json:"result"`
}
