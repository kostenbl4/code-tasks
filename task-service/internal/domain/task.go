package domain

import "github.com/google/uuid"

var (
	TaskResultError = "error"
	TaskResultOk    = "ok"
)

// Task - структура задачи
type Task struct {
	Translator string    `json:"translator"`
	Code       string    `json:"code"`
	ID         uuid.UUID `json:"task_id"`
	UserID     int64     `json:"-"`
	Status     string    `json:"status"`
	Result     string    `json:"result"`
	Stdout     string    `json:"stdout"`
	Stderr     string    `json:"stderr"`
}
