package domain

import "github.com/google/uuid"

// Task - структура задачи
type Task struct {
	UUID   uuid.UUID `json:"uuid"`
	Status string    `json:"status"`
	Result []byte    `json:"result"`
}
