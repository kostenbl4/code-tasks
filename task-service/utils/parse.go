package utils

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func ParseUUID(r *http.Request, param string) (uuid.UUID, error) {
	u := chi.URLParam(r, param)

	if u == "" {
		return uuid.UUID{}, fmt.Errorf("missing id")
	}

	id, err := uuid.Parse(u)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid id: %v", err)
	}
	return id, nil
}
