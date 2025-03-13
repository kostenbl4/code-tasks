package utils

import (
	"encoding/json"
	"net/http"
)

// ReadJSON - чтенеие json запроса
func ReadJSON(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	// decoder.DisallowUnknownFields() // для пропуска данных только с правильными полями
	return decoder.Decode(data)
}

// WriteJSON - запись json ответа
func WriteJSON(w http.ResponseWriter, data any, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

